package eod

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/lucasb-eyer/go-colorful"
)

const newText = "🆕"

var datafile *os.File
var createLock = &sync.Mutex{}

func (b *EoD) elemCreate(name string, parents []string, creator string, controversial string, guild string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}

	data := elems2txt(parents)
	_, res := dat.GetCombo(data)
	if res.Exists {
		return
	}

	_, res = dat.GetElement(name)
	text := "Combination"

	createLock.Lock()
	tx, err := b.db.BeginTx(context.Background(), nil)
	if err != nil {
		_ = tx.Rollback()
		createLock.Unlock()
		return
	}

	var postTxt string
	if !res.Exists {
		diff := -1
		compl := -1
		areUnique := false
		parColors := make([]int, len(parents))
		for j, val := range parents {
			elem, _ := dat.GetElement(val)
			if elem.Difficulty > diff {
				diff = elem.Difficulty
			}
			if elem.Complexity > compl {
				compl = elem.Complexity
			}
			if !strings.EqualFold(parents[0], val) {
				areUnique = true
			}
			parColors[j] = elem.Color
		}
		compl++
		if areUnique {
			diff++
		}
		col, err := mixColors(parColors)
		if err != nil {
			log.SetOutput(datafile)
			log.Println(err)
			tx.Rollback()
			createLock.Unlock()
			return
		}
		size, suc, msg := trees.ElemCreateSize(parents, dat)
		if !suc {
			log.SetOutput(datafile)
			log.Println(msg)
			tx.Rollback()
			createLock.Unlock()
			return
		}
		elem := types.Element{
			ID:         len(dat.Elements) + 1,
			Name:       name,
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  time.Now(),
			Parents:    parents,
			Complexity: compl,
			Difficulty: diff,
			Color:      col,
			TreeSize:   size,
		}
		postTxt = " - Element **#" + strconv.Itoa(elem.ID) + "**"
		dat.SetElement(elem)

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Color, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(parents), elem.Complexity, elem.Difficulty, 0, elem.TreeSize)
		if err != nil {
			dat.DeleteElement(elem.Name)

			log.SetOutput(datafile)
			log.Println(err)
			_ = tx.Rollback()
			createLock.Unlock()
			return
		}
		text = "Element"
	} else {
		el, res := dat.GetElement(name)
		if !res.Exists {
			log.SetOutput(datafile)
			log.Println("Doesn't exist")

			_ = tx.Rollback()
			createLock.Unlock()
			return
		}
		name = el.Name

		id := len(dat.Combos)
		if err == nil {
			postTxt = " - Combination **#" + strconv.Itoa(id) + "**"
		}
	}

	_, err = tx.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ? )", guild, data, name)
	if err != nil {
		dat.DeleteElement(name)
		createLock.Unlock()

		log.SetOutput(datafile)
		log.Println(err)
		_ = tx.Rollback()
		return
	}
	dat.AddComb(data, name)

	params := make(map[string]types.Empty)
	for _, val := range parents {
		params[val] = types.Empty{}
	}
	for k := range params {
		query := "UPDATE eod_elements SET usedin=usedin+1 WHERE name LIKE ? AND guild LIKE ?"
		if util.IsASCII(k) {
			query = "UPDATE eod_elements SET usedin=usedin+1 WHERE CONVERT(name USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND CONVERT(guild USING utf8mb4) LIKE CONVERT(? USING utf8mb4)"
		}
		if util.IsWildcard(k) {
			query = strings.ReplaceAll(query, " LIKE ", "=")
		}
		_, err = tx.Exec(query, k, guild)
		if err != nil {
			dat.DeleteElement(name)
			createLock.Unlock()

			log.SetOutput(datafile)
			log.Println(err)
			_ = tx.Rollback()
			return
		}

		el, res := dat.GetElement(k)
		if res.Exists {
			el.UsedIn++
			dat.SetElement(el)
		}
	}

	txt := newText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt + controversial

	_, _ = b.dg.ChannelMessageSend(dat.NewsChannel, txt)

	err = tx.Commit()
	if err != nil {
		dat.DeleteElement(name)
		createLock.Unlock()

		log.SetOutput(datafile)
		log.Println(err)
		return
	}

	createLock.Unlock()

	// Add Element to Inv
	inv, _ := dat.GetInv(creator, true)
	inv.Add(name)
	dat.SetInv(creator, inv)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()
	b.saveInv(guild, creator, true)

	err = b.autocategorize(name, guild)
	if err != nil {
		log.SetOutput(datafile)
		log.Println(err)
		return
	}
}

func mixColors(colors []int) (int, error) {
	cls := make([]colorful.Color, len(colors))
	var err error
	for i, color := range colors {
		hex := strconv.FormatInt(int64(color), 16)
		cls[i], err = colorful.Hex(hex)
		if err != nil {
			return 0, err
		}
	}

	var h, s, v float64
	for _, val := range cls {
		hv, sv, vv := val.Hsv()
		h += hv
		s += sv
		v += vv
	}
	length := float64(len(colors))
	h /= length
	s /= length
	v /= length

	out := colorful.Hsv(h, s, v)
	outv, err := strconv.ParseInt(out.Hex()[1:], 16, 64)
	return int(outv), err
}
