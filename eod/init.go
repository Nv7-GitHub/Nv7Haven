package eod

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/basecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/treecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/schollz/progressbar/v3"
)

func (b *EoD) init() {
	res, err := b.db.Query("SELECT * FROM eod_serverdata WHERE 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	var guild string
	var kind types.ServerDataType
	var value1 string
	var intval int
	for res.Next() {
		err = res.Scan(&guild, &kind, &value1, &intval)
		if err != nil {
			panic(err)
		}

		switch kind {
		case types.NewsChannel:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			dat.NewsChannel = value1
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.PlayChannel:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			if dat.PlayChannels == nil {
				dat.PlayChannels = make(map[string]types.Empty)
			}
			dat.PlayChannels[value1] = types.Empty{}
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.VotingChannel:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			dat.VotingChannel = value1
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.VoteCount:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			dat.VoteCount = intval
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.PollCount:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			dat.PollCount = intval
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.ModRole:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			dat.ModRole = value1
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()

		case types.UserColor:
			//lock.RLock()
			dat, exists := b.dat[guild]
			//lock.RUnlock()
			if !exists {
				dat = types.NewServerData()
			}
			if dat.UserColors == nil {
				dat.UserColors = make(map[string]int)
			}
			dat.UserColors[value1] = intval
			//lock.Lock()
			b.dat[guild] = dat
			//lock.Unlock()
		}
	}

	var cnt int
	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_elements").Scan(&cnt)
	if err != nil {
		panic(err)
	}

	bar := progressbar.New(cnt)

	elems, err := b.db.Query("SELECT name, image, color, guild, comment, creator, createdon, parents, complexity, difficulty, usedin, treesize FROM `eod_elements` ORDER BY createdon ASC")
	if err != nil {
		panic(err)
	}
	defer elems.Close()
	elem := types.Element{}
	var createdon int64
	var parentDat string
	for elems.Next() {
		err = elems.Scan(&elem.Name, &elem.Image, &elem.Color, &elem.Guild, &elem.Comment, &elem.Creator, &createdon, &parentDat, &elem.Complexity, &elem.Difficulty, &elem.UsedIn, &elem.TreeSize)
		if err != nil {
			return
		}
		elem.CreatedOn = time.Unix(createdon, 0)

		if len(parentDat) == 0 {
			elem.Parents = make([]string, 0)
		} else {
			elem.Parents = strings.Split(parentDat, "+")
		}

		//lock.RLock()
		dat, exists := b.dat[elem.Guild]
		//lock.RUnlock()
		if !exists {
			dat = types.NewServerData()
		}

		elem.ID = len(dat.Elements) + 1
		dat.Elements[strings.ToLower(elem.Name)] = elem
		//lock.Lock()
		b.dat[elem.Guild] = dat
		//lock.Unlock()

		bar.Add(1)
	}
	bar.Finish()

	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_combos").Scan(&cnt)
	if err != nil {
		panic(err)
	}

	bar = progressbar.New(cnt)

	combs, err := b.db.Query("SELECT * FROM `eod_combos`")
	if err != nil {
		panic(err)
	}
	defer combs.Close()
	var elemsVal string
	var elem3 string
	for combs.Next() {
		err = combs.Scan(&guild, &elemsVal, &elem3)
		if err != nil {
			return
		}
		//lock.RLock()
		dat, exists := b.dat[guild]
		//lock.RUnlock()
		if !exists {
			dat = types.NewServerData()
		}
		dat.Combos[elemsVal] = elem3
		//lock.Lock()
		b.dat[guild] = dat
		//lock.Unlock()

		bar.Add(1)
	}
	bar.Finish()

	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_elements").Scan(&cnt)
	if err != nil {
		panic(err)
	}

	bar = progressbar.New(cnt)

	invs, err := b.db.Query("SELECT guild, user, inv, made FROM eod_inv WHERE 1")
	if err != nil {
		panic(err)
	}
	defer invs.Close()
	var invDat string
	var user string
	var inv map[string]types.Empty
	var madecnt int
	for invs.Next() {
		inv = make(map[string]types.Empty)
		err = invs.Scan(&guild, &user, &invDat, &madecnt)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(invDat), &inv)
		if err != nil {
			panic(err)
		}
		//lock.RLock()
		dat, exists := b.dat[guild]
		//lock.RUnlock()
		if !exists {
			dat = types.NewServerData()
		}
		dat.Inventories[user] = types.Inventory{Elements: inv, MadeCnt: madecnt, User: user}
		//lock.Lock()
		b.dat[guild] = dat
		//lock.Unlock()

		bar.Add(1)
	}
	bar.Finish()

	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_categories").Scan(&cnt)
	if err != nil {
		panic(err)
	}
	bar = progressbar.New(cnt)

	cats, err := b.db.Query("SELECT * FROM eod_categories")
	if err != nil {
		panic(err)
	}
	defer cats.Close()
	var elemDat string
	cat := types.Category{}
	for cats.Next() {
		err = cats.Scan(&guild, &cat.Name, &elemDat, &cat.Image, &cat.Color)
		if err != nil {
			return
		}

		cat.Guild = guild

		//lock.RLock()
		dat, exists := b.dat[guild]
		//lock.RUnlock()
		if !exists {
			dat = types.NewServerData()
		}

		cat.Elements = make(map[string]types.Empty)
		err := json.Unmarshal([]byte(elemDat), &cat.Elements)
		if err != nil {
			panic(err)
		}

		dat.Categories[strings.ToLower(cat.Name)] = cat
		//lock.Lock()
		b.dat[guild] = dat
		//lock.Unlock()

		bar.Add(1)
	}

	bar.Finish()

	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_polls").Scan(&cnt)
	if err != nil {
		panic(err)
	}
	bar = progressbar.New(cnt)

	// Initialize subsystems
	logs.InitEoDLogs()
	b.base = base.NewBase(b.db, b.dat, b.dg, lock)
	b.treecmds = treecmds.NewTreeCmds(b.dat, b.dg, b.base, lock)
	b.polls = polls.NewPolls(b.dat, b.dg, b.db, b.base, lock)
	b.basecmds = basecmds.NewBaseCmds(b.dat, b.base, b.dg, b.db, lock)
	b.categories = categories.NewCategories(b.dat, b.base, b.dg, b.polls, lock)
	b.elements = elements.NewElements(b.dat, lock, b.polls, b.db, b.base, b.dg)

	polls, err := b.db.Query("SELECT * FROM eod_polls")
	if err != nil {
		panic(err)
	}
	defer polls.Close()
	var po types.Poll
	for polls.Next() {
		var jsondat string
		po.Data = nil
		err = polls.Scan(&po.Guild, &po.Channel, &po.Message, &po.Kind, &po.Value1, &po.Value2, &po.Value3, &po.Value4, &jsondat)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(jsondat), &po.Data)
		if err != nil {
			panic(err)
		}

		_, err = b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", po.Guild, po.Channel, po.Message)
		if err != nil {
			panic(err)
		}

		ups, err := b.dg.MessageReactions(po.Channel, po.Message, types.UpArrow, 100, "", "")
		if err != nil {
			panic(err)
		}
		po.Upvotes = len(ups) - 1

		downs, err := b.dg.MessageReactions(po.Channel, po.Message, types.DownArrow, 100, "", "")
		if err != nil {
			panic(err)
		}
		po.Downvotes = len(downs) - 1

		b.dat[po.Guild], _ = b.polls.CheckReactions(b.dat[po.Guild], po, downs[len(downs)-1].ID)

		bar.Add(1)
	}

	bar.Finish()

	b.initHandlers()
	b.start()

	// Start stats saving
	go func() {
		b.basecmds.SaveStats()
		for {
			time.Sleep(time.Minute * 30)
			b.basecmds.SaveStats()
		}
	}()

	// Recalc autocats?
	if types.RecalcAutocats {
		for id, gld := range b.dat {
			for elem := range gld.Elements {
				b.polls.Autocategorize(elem, id)
			}
		}
	}
}
