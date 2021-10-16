package base

import (
	"fmt"
	"math"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Base) InvPageGetter(p types.PageSwitcher) (string, int, int, error) {
	length := int(math.Floor(float64(len(p.Items)-1) / float64(p.PageLength)))
	if p.PageLength*p.Page > (len(p.Items) - 1) { // If you go past the max page, bring you to 0
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

func (b *Base) LbPageGetter(p types.PageSwitcher) (string, int, int, error) {
	length := int(math.Floor(float64(len(p.Users)-1) / float64(p.PageLength)))
	if p.PageLength*p.Page > (len(p.Users) - 1) { // If you go past the max page, bring you to 0
		return "", 0, length, nil
	}

	if p.Page < 0 { // If you go before the first page, go to the last one
		return "", length, length, nil
	}

	txt := &strings.Builder{}

	// Get the range
	first := p.PageLength * p.Page
	users := p.Users[first:]
	if len(users) > p.PageLength {
		users = users[:p.PageLength]
	}
	max := first + len(users)

	// Create the text
	containsYou := false
	for i := first; i < max; i++ {
		ft := "%d. <@%s> - %d\n"
		if p.Users[i] == p.User {
			ft = "%d. <@%s> *You* - %d\n"
			containsYou = true
		}
		fmt.Fprintf(txt, ft, i+1, p.Users[i], p.Cnts[i])
	}

	if !containsYou {
		fmt.Fprintf(txt, "\n%d. <@%s> *You* - %d", p.UserPos+1, p.Users[p.UserPos], p.Cnts[p.UserPos])
	}
	return txt.String(), p.Page, length, nil
}

func (b *Base) SearchPageGetter(p types.PageSwitcher) (string, int, int, error) {
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
	res, err := b.db.Query("SELECT name FROM eod_elements WHERE guild=? AND LOWER(name) LIKE ? ORDER BY name LIMIT ? OFFSET ?", p.Guild, wild, p.PageLength, p.Page*p.PageLength)
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
