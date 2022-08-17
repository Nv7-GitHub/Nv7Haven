package elements

import (
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Elements) InvCmd(user string, m types.Msg, rsp types.Rsp, sorter string, filter string, postfix bool, defaul bool) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	inv := db.GetInv(user)
	type invItem struct {
		name string
		id   int
	}
	items := make([]invItem, len(inv.Elements))
	i := 0
	db.RLock()
	for k := range inv.Elements {
		el, _ := db.GetElement(k, true)
		items[i] = invItem{el.Name, el.ID}
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]invItem, len(items))
		for _, val := range items {
			creator := ""
			elem, res := db.GetElement(val.id, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = invItem{elem.Name, elem.ID}
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}
	if defaul && m.Author.ID != user {
		postfix = true
	}
	eodsort.Sort(items, len(items), func(index int) int {
		return items[index].id
	}, func(index int) string {
		return items[index].name
	}, func(index int, val string) {
		items[index].name = val
	}, sorter, m.Author.ID, db, postfix)

	text := make([]string, len(items))
	for i, v := range items {
		text[i] = v.name
	}
	db.RUnlock()

	name := m.Author.Username
	if m.Author.ID != user {
		u, err := b.dg.User(user)
		if rsp.Error(err) {
			return
		}
		name = u.Username
	}
	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("UserInventory", map[string]any{
			"Username": name,
			"Count":    len(items),
			"Percent":  util.FormatFloat(float32(len(items))/float32(len(db.Elements))*100, 2),
		}),
		PageGetter: b.base.InvPageGetter,
		Items:      text,
	}, m, rsp)
}

type invInfo struct {
	User          string
	ElementCnt    int
	MadeCnt       int
	SignedCnt     int
	ImagedCnt     int
	ColoredCnt    int
	CatImagedCnt  int
	CatColoredCnt int
	CatSignedCnt  int
	UsedCnt       int
}

func (b *Elements) LbCmd(m types.Msg, rsp types.Rsp, sorter string, user string, category string) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	db.GetInv(user) // Make user exist

	var cat []types.Element
	if category != "" {
		catV, res := db.GetCat(category)
		if !res.Exists {
			catV, res := db.GetVCat(category)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			}
			els, res := b.base.CalcVCat(catV, db, true)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			}
			category = catV.Name
			cat = make([]types.Element, 0, len(els))
			db.RLock()
			for k := range els {
				el, res := db.GetElement(k, true)
				if res.Exists {
					cat = append(cat, el)
				}
			}
			db.RUnlock()
		} else {
			category = catV.Name
			cat = make([]types.Element, 0, len(catV.Elements))
			db.RLock()
			for k := range catV.Elements {
				el, res := db.GetElement(k, true)
				if res.Exists {
					cat = append(cat, el)
				}
			}
			db.RUnlock()
		}
	}

	// Sort invs
	invs := make([]invInfo, len(db.Invs()))
	i := 0
	for _, v := range db.Invs() {
		if cat != nil {
			v := invInfo{User: v.User}
			switch sorter {
			case "made":
				for _, el := range cat {
					if el.Creator == v.User {
						v.MadeCnt++
					}
				}

			case "signed":
				for _, el := range cat {
					if el.Commenter == v.User {
						v.SignedCnt++
					}
				}

			case "imaged":
				for _, el := range cat {
					if el.Imager == v.User {
						v.ImagedCnt++
					}
				}

			case "colored":
				for _, el := range cat {
					if el.Colorer == v.User {
						v.ColoredCnt++
					}
				}

			default:
				inv := db.GetInv(v.User)
				inv.Lock.RLock()
				for _, el := range cat {
					_, exists := inv.Elements[el.ID]
					if exists {
						v.ElementCnt++
					}
				}
				inv.Lock.RUnlock()
			}
			invs[i] = v
		} else {
			invs[i] = invInfo{
				User:          v.User,
				ElementCnt:    len(v.Elements),
				MadeCnt:       v.MadeCnt,
				SignedCnt:     v.SignedCnt,
				ImagedCnt:     v.ImagedCnt,
				ColoredCnt:    v.ColoredCnt,
				CatImagedCnt:  v.CatImagedCnt,
				CatColoredCnt: v.CatColoredCnt,
				CatSignedCnt:  v.CatSignedCnt,
				UsedCnt:       v.UsedCnt,
			}
		}
		i++
	}
	var sortFn func(a, b int) bool
	var titleID string
	switch sorter {
	case "made":
		sortFn = func(a, b int) bool {
			return invs[a].MadeCnt > invs[b].MadeCnt
		}
		titleID = "LbTitleMade"

	case "signed":
		sortFn = func(a, b int) bool {
			return invs[a].SignedCnt > invs[b].SignedCnt
		}
		titleID = "LbTitleSigned"

	case "imaged":
		sortFn = func(a, b int) bool {
			return invs[a].ImagedCnt > invs[b].ImagedCnt
		}
		titleID = "LbTitleImaged"

	case "colored":
		sortFn = func(a, b int) bool {
			return invs[a].ColoredCnt > invs[b].ColoredCnt
		}
		titleID = "LbTitleColored"

	case "catimaged":
		sortFn = func(a, b int) bool {
			return invs[a].CatImagedCnt > invs[b].CatImagedCnt
		}
		titleID = "LbTitleCatImaged"

	case "catcolored":
		sortFn = func(a, b int) bool {
			return invs[a].CatColoredCnt > invs[b].CatColoredCnt
		}
		titleID = "LbTitleCatColored"

	case "catsigned":
		sortFn = func(a, b int) bool {
			return invs[a].CatSignedCnt > invs[b].CatSignedCnt
		}
		titleID = "LbTitleCatSigned"

	case "used":
		sortFn = func(a, b int) bool {
			return invs[a].UsedCnt > invs[b].UsedCnt
		}
		titleID = "LbTitleUsed"

	default:
		sortFn = func(a, b int) bool {
			return invs[a].ElementCnt > invs[b].ElementCnt
		}
		titleID = "LbTitleElem"
	}
	sort.Slice(invs, sortFn)

	// Convert to right format
	users := make([]string, len(invs))
	cnts := make([]int, len(invs))
	userpos := 0
	for i, v := range invs {
		users[i] = v.User
		switch sorter {
		case "made":
			cnts[i] = v.MadeCnt

		case "signed":
			cnts[i] = v.SignedCnt

		case "imaged":
			cnts[i] = v.ImagedCnt

		case "colored":
			cnts[i] = v.ColoredCnt

		case "catimaged":
			cnts[i] = v.CatImagedCnt

		case "catcolored":
			cnts[i] = v.CatColoredCnt

		case "used":
			cnts[i] = v.UsedCnt

		case "catsigned":
			cnts[i] = v.CatSignedCnt

		default:
			cnts[i] = v.ElementCnt
		}
		if v.User == user {
			userpos = i
		}
	}

	title := db.Config.LangProperty(titleID, nil)
	if cat != nil {
		title += " (" + category + ")"
	}
	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      title,
		PageGetter: b.base.LbPageGetter,

		User:    user,
		Users:   users,
		UserPos: userpos,
		Cnts:    cnts,
	}, m, rsp)
}
