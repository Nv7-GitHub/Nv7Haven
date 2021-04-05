package eod

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *EoD) init() {
	// Debugging
	var err error
	datafile, err = os.OpenFile("eodlogs.txt", os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	b.initInfoChoices()
	for _, v := range commands {
		go func(val *discordgo.ApplicationCommand) {
			if val.Name == "elemsort" {
				val.Options[0].Choices = infoChoices
			}
			_, err := b.dg.ApplicationCommandCreate(clientID, "", val) // 819077688371314718 for testing
			if err != nil {
				panic(err)
			}
		}(v)
	}
	b.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rsp := b.newRespSlash(i)
		if (i.Data.Name != "suggest") && (i.Data.Name != "mark") && (i.Data.Name != "image") && (i.Data.Name != "inv") && (i.Data.Name != "lb") && (i.Data.Name != "addcat") && (i.Data.Name != "cat") && (i.Data.Name != "hint") && (i.Data.Name != "stats") && (i.Data.Name != "idea") && (i.Data.Name != "about") && (i.Data.Name != "path") && (i.Data.Name != "get") && (i.Data.Name != "rmcat") {
			isMod, err := b.isMod(i.Member.User.ID, i.GuildID, bot.newMsgSlash(i))
			if rsp.Error(err) {
				return
			}
			if !isMod {
				rsp.ErrorMessage("You need to have permission `Administrator`!")
				return
			}
		}
		if i.Data.Name == "path" {
			isMod, err := b.isMod(i.Member.User.ID, i.GuildID, bot.newMsgSlash(i))
			if rsp.Error(err) {
				return
			}
			if !isMod {
				lock.RLock()
				dat, exists := b.dat[i.GuildID]
				lock.RUnlock()
				if !exists {
					rsp.ErrorMessage("You need to have permission `Administrator`!")
					return
				}
				inv, exists := dat.invCache[i.Member.User.ID]
				if !exists {
					rsp.ErrorMessage("You need to have permission `Administrator`!")
					return
				}
				_, exists = inv[strings.ToLower(i.Data.Options[0].StringValue())]
				if !exists {
					rsp.ErrorMessage("You don't have that element!")
					return
				}
			}
		}
		if h, ok := commandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})
	b.dg.AddHandler(b.cmdHandler)
	b.dg.AddHandler(b.reactionHandler)
	b.dg.AddHandler(b.unReactionHandler)
	b.dg.AddHandler(b.pageSwitchHandler)
	b.dg.AddHandler(func(s *discordgo.Session, i *discordgo.Disconnect) {
		log.Println("Disconnected!")
	})

	res, err := b.db.Query("SELECT * FROM eod_serverdata WHERE 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	var guild string
	var kind serverDataType
	var value1 string
	var intval int
	for res.Next() {
		err = res.Scan(&guild, &kind, &value1, &intval)
		if err != nil {
			panic(err)
		}

		switch kind {
		case newsChannel:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			dat.newsChannel = value1
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()

		case playChannel:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			if dat.playChannels == nil {
				dat.playChannels = make(map[string]empty)
			}
			dat.playChannels[value1] = empty{}
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()

		case votingChannel:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			dat.votingChannel = value1
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()

		case voteCount:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			dat.voteCount = intval
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()

		case pollCount:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			dat.pollCount = intval
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()

		case modRole:
			lock.RLock()
			dat, exists := b.dat[guild]
			lock.RUnlock()
			if !exists {
				dat = serverData{}
			}
			dat.modRole = value1
			lock.Lock()
			b.dat[guild] = dat
			lock.Unlock()
		}
	}

	elems, err := b.db.Query("SELECT * FROM eod_elements WHERE 1")
	if err != nil {
		panic(err)
	}
	defer elems.Close()
	elem := element{}
	var createdon int64
	var catDat string
	var parentDat string
	for elems.Next() {
		err = elems.Scan(&elem.Name, &catDat, &elem.Image, &elem.Guild, &elem.Comment, &elem.Creator, &createdon, &parentDat, &elem.Complexity, &elem.Difficulty, &elem.UsedIn)
		if err != nil {
			return
		}
		elem.Categories = make(map[string]empty)
		err = json.Unmarshal([]byte(catDat), &elem.Categories)
		if err != nil {
			panic(err)
		}
		elem.CreatedOn = time.Unix(createdon, 0)
		parentMap := make(map[string]empty)
		err = json.Unmarshal([]byte(parentDat), &parentMap)
		if err != nil {
			return
		}
		parents := make([]string, len(parentMap))
		i := 0
		for k := range parentMap {
			parents[i] = k
			i++
		}
		elem.Parents = parents

		lock.RLock()
		dat := b.dat[elem.Guild]
		lock.RUnlock()
		if dat.elemCache == nil {
			dat.elemCache = make(map[string]element)
		}
		dat.elemCache[strings.ToLower(elem.Name)] = elem
		lock.Lock()
		b.dat[elem.Guild] = dat
		lock.Unlock()
	}

	invs, err := b.db.Query("SELECT guild, user, inv FROM eod_inv WHERE 1")
	if err != nil {
		panic(err)
	}
	defer invs.Close()
	var invDat string
	var user string
	var inv map[string]empty
	for invs.Next() {
		inv = make(map[string]empty)
		err = invs.Scan(&guild, &user, &invDat)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(invDat), &inv)
		if err != nil {
			panic(err)
		}
		lock.RLock()
		dat := b.dat[guild]
		lock.RUnlock()
		if dat.invCache == nil {
			dat.invCache = make(map[string]map[string]empty)
		}
		dat.invCache[user] = inv
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
	}

	polls, err := b.db.Query("SELECT * FROM eod_polls WHERE 1")
	if err != nil {
		panic(err)
	}
	defer polls.Close()
	var po poll
	for polls.Next() {
		var jsondat string
		err = polls.Scan(&guild, &po.Channel, &po.Message, &po.Kind, &po.Value1, &po.Value2, &po.Value3, &po.Value4, &jsondat)
		if err != nil {
			panic(err)
		}
		po.Guild = guild
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
	}
}
