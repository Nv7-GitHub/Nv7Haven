package eod

import (
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// element, guild, guild, element, guild, element - returns: made by x combos, used in x combos, found by x people
const elemInfoDataCount = `SELECT a.cnt, b.cnt, c.cnt FROM (SELECT COUNT(1) AS cnt FROM eod_combos WHERE elem3=? AND guild=?) a, (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)) b, (SELECT e.rw AS cnt FROM (SELECT ROW_NUMBER() OVER (ORDER BY createdon ASC) AS rw, name FROM eod_elements WHERE guild=?) e WHERE e.name=?) c`

var infoChoices []*discordgo.ApplicationCommandOptionChoice
var infoQuerys = map[string]string{
	"Name":         "SELECT name FROM eod_elements WHERE guild=? ORDER BY name %s LIMIT ? OFFSET ?",
	"Date Created": "SELECT name FROM eod_elements WHERE guild=? ORDER BY createdon %s LIMIT ? OFFSET ?",
	"Complexity":   "SELECT name FROM eod_elements WHERE guild=? ORDER BY complexity %s LIMIT ? OFFSET ?",
	"Difficulty":   "SELECT name FROM eod_elements WHERE guild=? ORDER BY difficulty %s LIMIT ? OFFSET ?",
	"Used In":      `SELECT name FROM eod_elements WHERE guild=? ORDER BY usedin %s LIMIT ? OFFSET ?`,
	"Made By":      "SELECT name FROM eod_elements WHERE guild=? ORDER BY (SELECT COUNT(1) AS cnt FROM eod_combos WHERE elem3=name AND guild=?) %s LIMIT ? OFFSET ?",
	"Found By":     `SELECT name FROM eod_elements WHERE guild=?  ORDER BY (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(name), '"')) IS NOT NULL)) %s LIMIT ? OFFSET ?`,
}

func (b *EoD) initInfoChoices() {
	infoChoices = make([]*discordgo.ApplicationCommandOptionChoice, len(infoQuerys))
	i := 0
	for k := range infoQuerys {
		infoChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  k,
			Value: k,
		}
		i++
	}
}

func (b *EoD) sortPageGetter(p pageSwitcher) (string, int, int, error) {
	length := int(math.Floor(float64(p.Length-1) / float64(pageLength)))
	if pageLength*p.Page > (p.Length - 1) {
		return "", 0, length, nil
	}
	if p.Page < 0 {
		return "", length, length, nil
	}
	var res *sql.Rows
	var err error
	cnt := strings.Count(p.Query, "?")
	if cnt == 3 {
		res, err = b.db.Query(p.Query, p.Guild, pageLength, p.Page*pageLength)
		if err != nil {
			return "", length, length, err
		}
	} else {
		res, err = b.db.Query(p.Query, p.Guild, p.Guild, pageLength, p.Page*pageLength)
		if err != nil {
			return "", length, length, err
		}
	}
	defer res.Close()
	out := ""
	var name string
	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			return "", length, length, err
		}
		out += name + "\n"
	}
	return out, p.Page, length, nil
}

func (b *EoD) sortCmd(query string, order bool, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	quer, exists := infoQuerys[query]
	if !exists {
		rsp.ErrorMessage("Invalid query type!")
		return
	}
	ord := "DESC"
	if order {
		ord = "ASC"
	}
	quer = fmt.Sprintf(quer, ord)
	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchElemSort,
		Title:      "Element Sort",
		PageGetter: b.sortPageGetter,
		Query:      quer,
		Length:     len(dat.elemCache),
	}, m, rsp)
}

func (b *EoD) infoCmd(elem string, isNumber bool, number int, m msg, rsp rsp) {
	rsp.Acknowledge()

	if len(elem) == 0 {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}
	if isNumber {
		row := b.db.QueryRow(`SELECT e.name AS cnt FROM (SELECT ROW_NUMBER() OVER (ORDER BY createdon ASC) AS rw, name FROM eod_elements WHERE guild=?) e WHERE e.rw=?`, m.GuildID, number)
		row.Scan(&elem)
	}
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		if strings.Contains(elem, "?") {
			isValid := false
			for _, letter := range elem {
				if letter != '?' {
					isValid = true
					break
				}
			}
			if !isValid {
				rsp.ErrorMessage("Invalid letter!")
				return
			}
		}
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	has := ""
	exists = false
	if dat.invCache != nil {
		_, exists = dat.invCache[m.Author.ID]
		if exists {
			_, exists = dat.invCache[m.Author.ID][strings.ToLower(el.Name)]
		}
	}
	if !exists {
		has = "don't "
	}

	row := b.db.QueryRow(elemInfoDataCount, el.Name, el.Guild, el.Guild, el.Name, el.Guild, el.Name)
	var madeby int
	var foundby int
	var id int
	err := row.Scan(&madeby, &foundby, &id)
	if rsp.Error(err) {
		return
	}

	usedbysuff := "s"
	if el.UsedIn == 1 {
		usedbysuff = ""
	}
	madebysuff := "s"
	if madeby == 1 {
		madebysuff = ""
	}
	foundbysuff := "s"
	if foundby == 1 {
		foundbysuff = ""
	}

	rsp.Embed(&discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Element **#%d**\nCreated by <@%s>\nCreated on %s\nUsed in %d combo%s\nMade with %d combo%s\nFound by %d player%s\nComplexity: %d\nDifficulty: %d\n<@%s> **You %shave this.**\n\n%s", id, el.Creator, el.CreatedOn.Format("January 2, 2006, 3:04 PM"), el.UsedIn, usedbysuff, madeby, madebysuff, foundby, foundbysuff, el.Complexity, el.Difficulty, m.Author.ID, has, el.Comment),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
	})
}
