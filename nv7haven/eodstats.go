package nv7haven

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type eodStats struct {
	refreshTime time.Time
	Labels      []string `json:"labels"`

	Found []int `json:"found"`

	Elemcnt     []int `json:"elemcnt"`
	Categorized []int `json:"categorized"`
	Combcnt     []int `json:"combcnt"`

	Usercnt   []int `json:"usercnt"`
	Servercnt []int `json:"servercnt"`

	chart string
}

func (n *Nv7Haven) refreshStats() {
	res, err := n.sql.Query("SELECT * FROM eod_stats WHERE time > ? ORDER BY time ", n.eodStats.refreshTime.Unix())
	if err != nil {
		fmt.Println(err)
	}

	changed := false

	var tm int64
	var elemcnt, combcnt, usercnt, found, categorized, servercnt int
	for res.Next() {
		err = res.Scan(&tm, &elemcnt, &combcnt, &usercnt, &found, &categorized, &servercnt)
		if err != nil {
			fmt.Println(err)
		}

		n.eodStats.Labels = append(n.eodStats.Labels, time.Unix(tm, 0).Format("2006-01-02"))
		n.eodStats.Elemcnt = append(n.eodStats.Elemcnt, elemcnt)
		n.eodStats.Combcnt = append(n.eodStats.Combcnt, combcnt)
		n.eodStats.Usercnt = append(n.eodStats.Usercnt, usercnt)
		n.eodStats.Servercnt = append(n.eodStats.Servercnt, servercnt)
		n.eodStats.Found = append(n.eodStats.Found, found)
		n.eodStats.Categorized = append(n.eodStats.Categorized, categorized)

		if !changed {
			changed = true
			n.eodStats.refreshTime = time.Now()
		}
	}

	if changed {
		dat, err := json.Marshal(n.eodStats)
		if err != nil {
			fmt.Println(err)
		}
		n.eodStats.chart = string(dat)
	}
}

func (n *Nv7Haven) getEodStats(c *fiber.Ctx) error {
	go n.refreshStats()
	c.WriteString(n.eodStats.chart)
	return nil
}
