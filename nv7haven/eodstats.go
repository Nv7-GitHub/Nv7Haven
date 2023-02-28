package nv7haven

import (
	"encoding/json"
	"fmt"
	"sort"
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

	CommandCounts []CommandCount `json:"commandcounts"`

	chart string
}

type CommandCount struct {
	Counts     map[string]int `json:"counts"`
	Time       int64          `json:"time"`
	TimeString string         `json:"timestring"`
}

func (n *Nv7Haven) refreshStats() {
	res, err := n.pgdb.Query("SELECT * FROM stats WHERE time > $1 ORDER BY time", n.eodStats.refreshTime)
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
		}
	}
	res.Close()

	res, err = n.pgdb.Query("SELECT * FROM stats_commands WHERE time > $1 ORDER BY time", n.eodStats.refreshTime)
	if err != nil {
		fmt.Println(err)
	}

	var name string
	var cnt int
	for res.Next() {
		err = res.Scan(&tm, &name, &cnt)
		if err != nil {
			fmt.Println(err)
		}

		// Find tm
		found := false
		for i, t := range n.eodStats.CommandCounts {
			if t.Time == tm {
				t.Counts[name] = cnt
				n.eodStats.CommandCounts[i] = t
				found = true
				break
			}
		}
		if !found {
			n.eodStats.CommandCounts = append(n.eodStats.CommandCounts, CommandCount{
				Counts:     map[string]int{name: cnt},
				Time:       tm,
				TimeString: time.Unix(tm, 0).Format("2006-01-02"),
			})
		}

		if !changed {
			changed = true
		}
	}
	res.Close()

	sort.Slice(n.eodStats.CommandCounts, func(i, j int) bool {
		return n.eodStats.CommandCounts[i].Time < n.eodStats.CommandCounts[j].Time
	})

	if changed {
		n.eodStats.refreshTime = time.Now()

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
