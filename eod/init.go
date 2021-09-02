package eod

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	//elems, err := b.db.Query("SELECT * FROM eod_elements ORDER BY createdon ASC") // Do after nov 21

	var cnt int
	err = b.db.QueryRow("SELECT COUNT(1) FROM eod_elements").Scan(&cnt)
	if err != nil {
		panic(err)
	}

	bar := progressbar.New(cnt)

	elems, err := b.db.Query("SELECT name, image, guild, comment, creator, createdon, parents, complexity, difficulty, usedin FROM `eod_elements` ORDER BY (CASE WHEN createdon=1637536881 THEN 1605988759 ELSE createdon END) ASC")
	if err != nil {
		panic(err)
	}
	defer elems.Close()
	elem := types.Element{}
	var createdon int64
	var parentDat string
	for elems.Next() {
		err = elems.Scan(&elem.Name, &elem.Image, &elem.Guild, &elem.Comment, &elem.Creator, &createdon, &parentDat, &elem.Complexity, &elem.Difficulty, &elem.UsedIn)
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
		dat := b.dat[elem.Guild]
		//lock.RUnlock()
		if dat.Elements == nil {
			dat.Elements = make(map[string]types.Element)
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
		dat := b.dat[guild]
		//lock.RUnlock()
		if dat.Combos == nil {
			dat.Combos = make(map[string]string)
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

	invs, err := b.db.Query("SELECT guild, user, inv FROM eod_inv WHERE 1")
	if err != nil {
		panic(err)
	}
	defer invs.Close()
	var invDat string
	var user string
	var inv map[string]types.Empty
	for invs.Next() {
		inv = make(map[string]types.Empty)
		err = invs.Scan(&guild, &user, &invDat)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(invDat), &inv)
		if err != nil {
			panic(err)
		}
		//lock.RLock()
		dat := b.dat[guild]
		//lock.RUnlock()
		if dat.Inventories == nil {
			dat.Inventories = make(map[string]types.Container)
		}
		dat.Inventories[user] = inv
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

		b.dg.ChannelMessageDelete(po.Channel, po.Message)
		err = b.createPoll(po)
		if err != nil {
			fmt.Println(err)
		}
		bar.Add(1)
	}

	bar.Finish()

	//lock.RLock()
	for k, dat := range b.dat {
		hasChanged := false
		if dat.Inventories == nil {
			dat.Inventories = make(map[string]types.Container)
			hasChanged = true
		}
		if hasChanged {
			//lock.RUnlock()
			//lock.Lock()
			b.dat[k] = dat
			//lock.Unlock()
			//lock.RLock()
		}
	}
	//lock.RUnlock()

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
		err = cats.Scan(&guild, &cat.Name, &elemDat, &cat.Image)
		if err != nil {
			return
		}

		cat.Guild = guild

		//lock.RLock()
		dat := b.dat[guild]
		//lock.RUnlock()
		if dat.Categories == nil {
			dat.Categories = make(map[string]types.Category)
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

	b.initHandlers()

	// Start stats saving
	go func() {
		b.saveStats()
		for {
			time.Sleep(time.Minute * 30)
			b.saveStats()
		}
	}()
}
