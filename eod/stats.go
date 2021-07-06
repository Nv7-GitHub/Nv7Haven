package eod

import (
	"fmt"
	"log"
	"time"
)

func (b *EoD) statsCmd(m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	gd, err := b.dg.State.Guild(m.GuildID)
	if rsp.Error(err) {
		return
	}

	var cnt int
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=?", m.GuildID)
	err = row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	found := 0
	dat.lock.RLock()
	for _, val := range dat.invCache {
		found += len(val)
	}

	categorized := 0
	for _, val := range dat.catCache {
		categorized += len(val.Elements)
	}
	dat.lock.RUnlock()

	dat.lock.RLock()
	rsp.Message(fmt.Sprintf("Element Count: **%s**\nCombination Count: **%s**\nMember Count: **%s**\nElements Found: **%s**\nElements Categorized: **%s**", formatInt(len(dat.elemCache)), formatInt(cnt), formatInt(gd.MemberCount), formatInt(found), formatInt(categorized)))
	dat.lock.RUnlock()
}

// takes time, found, categorized
var saveStatsQuery = `INSERT INTO eod_stats VALUES (?, (SELECT COUNT(1) FROM eod_elements), (SELECT COUNT(1) FROM eod_combos), (SELECT COUNT(DISTINCT user) FROM eod_inv), ?, ?, (SELECT COUNT(DISTINCT guild) FROM eod_serverdata))`

func (b *EoD) saveStats() {
	var lastTime int64
	err := b.db.QueryRow("SELECT time FROM eod_stats ORDER BY time DESC LIMIT 1").Scan(&lastTime)
	if err != nil {
		log.Println(err)
	}

	if time.Since(time.Unix(lastTime, 0)).Hours() > 24 {
		lock.RLock()
		categorized := 0
		found := 0
		for _, dat := range b.dat {
			for _, val := range dat.catCache {
				categorized += len(val.Elements)
			}
			for _, val := range dat.invCache {
				found += len(val)
			}
		}
		lock.RUnlock()

		_, err = b.db.Exec(saveStatsQuery, time.Now().Unix(), found, categorized)
		if err != nil {
			log.Println(err)
		}
	}
}
