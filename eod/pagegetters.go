package eod

import (
	"fmt"
	"math"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

const ldbQuery = `
SELECT rw, ` + "user" + `, ` + "%s" + `
FROM (
    SELECT 
         ROW_NUMBER() OVER (ORDER BY ` + "%s" + ` DESC) AS rw,
         ` + "user" + `, ` + "%s" + `
    FROM eod_inv WHERE guild=?
) sub
WHERE sub.user=?
`

func (b *EoD) invPageGetter(p types.PageSwitcher) (string, int, int, error) {
	length := int(math.Floor(float64(len(p.Items)-1) / float64(p.PageLength)))
	if p.PageLength*p.Page > (len(p.Items) - 1) {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", length, length, nil
	}

	items := p.Items[p.PageLength*p.Page:]
	if len(items) > p.PageLength {
		items = items[:p.PageLength]
	}
	return strings.Join(items, "\n"), p.Page, length, nil
}

func (b *EoD) lbPageGetter(p types.PageSwitcher) (string, int, int, error) {
	cnt := b.db.QueryRow("SELECT COUNT(1) FROM eod_inv WHERE guild=?", p.Guild)
	pos := b.db.QueryRow(fmt.Sprintf(ldbQuery, p.Sort, p.Sort, p.Sort), p.Guild, p.User)
	var count int
	var ps int
	var u string
	var ucnt int
	err := pos.Scan(&ps, &u, &ucnt)
	if err != nil {
		return "", 0, 0, err
	}
	cnt.Scan(&count)
	length := int(math.Floor(float64(count-1) / float64(p.PageLength)))
	if err != nil {
		return "", 0, 0, err
	}
	if p.PageLength*p.Page > (count - 1) {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", length, length, nil
	}

	text := ""
	res, err := b.db.Query(fmt.Sprintf("SELECT %s, `user` FROM eod_inv WHERE guild=? ORDER BY %s DESC LIMIT ? OFFSET ?", p.Sort, p.Sort), p.Guild, p.PageLength, p.Page*p.PageLength)
	if err != nil {
		return "", 0, 0, err
	}
	defer res.Close()
	i := p.PageLength*p.Page + 1
	var user string
	var ct int
	for res.Next() {
		err = res.Scan(&ct, &user)
		if err != nil {
			return "", 0, 0, err
		}
		you := ""
		if user == u {
			you = " *You*"
		}
		text += fmt.Sprintf("%d. <@%s>%s - %d\n", i, user, you, ct)
		i++
	}
	if !((p.PageLength*p.Page <= ps) && (ps <= (p.Page+1)*p.PageLength)) {
		text += fmt.Sprintf("\n%d. <@%s> *You* - %d\n", ps, u, ucnt)
	}
	return text, p.Page, length, nil
}

func (b *EoD) searchPageGetter(p types.PageSwitcher) (string, int, int, error) {
	wild := "%" + util.EscapeElement(strings.ToLower(p.Search)) + "%"

	var count int
	err := b.db.QueryRow("SELECT COUNT(1) FROM eod_elements WHERE guild=? AND LOWER(name) LIKE ?", p.Guild, wild).Scan(&count)
	if err != nil {
		return "", 0, 0, err
	}

	length := int(math.Floor(float64(count-1) / float64(p.PageLength)))
	if p.PageLength*p.Page > (count - 1) {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", length, length, nil
	}

	text := ""
	res, err := b.db.Query("SELECT name FROM eod_elements WHERE guild=? AND LOWER(name) LIKE ? ORDER BY name ASC LIMIT ? OFFSET ?", p.Guild, wild, p.PageLength, p.Page*p.PageLength)
	if err != nil {
		return "", 0, 0, err
	}
	defer res.Close()
	i := p.PageLength*p.Page + 1
	var elem string
	for res.Next() {
		err = res.Scan(&elem)
		if err != nil {
			return "", 0, 0, err
		}
		text += elem + "\n"
		i++
	}
	return text, p.Page, length, nil
}
