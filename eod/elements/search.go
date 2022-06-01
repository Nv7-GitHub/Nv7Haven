package elements

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Elements) SearchCmd(search string, sort string, source string, opt string, regex bool, postfix bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	var list map[string]types.Empty
	switch source {
	case "elements":
		list = make(map[string]types.Empty, len(db.Elements))
		for _, el := range db.Elements {
			list[el.Name] = types.Empty{}
		}

	case "inventory":
		inv := db.GetInv(opt)

		list = make(map[string]types.Empty, len(inv.Elements))
		inv.Lock.RLock()
		db.RLock()
		for el := range inv.Elements {
			elem, res := db.GetElement(el, true)
			if !res.Exists {
				continue
			}
			list[elem.Name] = types.Empty{}
		}
		db.RUnlock()
		inv.Lock.RUnlock()

	case "category":
		var els map[int]types.Empty
		var lock *sync.RWMutex
		catv, res := db.GetCat(opt)
		if !res.Exists {
			vcat, res := db.GetVCat(opt)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			}
			els, res = b.base.CalcVCat(vcat, db)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			}
		} else {
			els = catv.Elements
			lock = catv.Lock
		}

		list = make(map[string]types.Empty, len(els))
		if lock != nil {
			lock.RLock()
		}
		db.RLock()
		for el := range els {
			elem, res := db.GetElement(el, true)
			if !res.Exists {
				continue
			}
			list[elem.Name] = types.Empty{}
		}
		db.RUnlock()
		if lock != nil {
			lock.RUnlock()
		}
	}

	type searchResult struct {
		name string
		id   int
	}
	results := make([]searchResult, 0)
	if regex {
		reg, err := regexp.Compile(search)
		if rsp.Error(err) {
			return
		}
		db.RLock()
		for el := range list {
			m := reg.Find([]byte(el))
			if m != nil {
				el, res := db.GetElementByName(string(m), true)
				if !res.Exists {
					continue
				}
				results = append(results, searchResult{name: el.Name, id: el.ID})
			}
		}
		db.RUnlock()
	} else {
		s := strings.ToLower(search)
		db.RLock()
		for el := range list {
			if strings.Contains(strings.ToLower(el), s) {
				elem, res := db.GetElementByName(el, true)
				if !res.Exists {
					continue
				}
				name := []rune(elem.Name)
				pos := strings.Index(strings.ToLower(el), s)
				fmt.Println(elem.Name, s, pos)
				name = append(name[:pos+len(s)], append([]rune("**"), name[pos+len(s):]...)...)
				name = append(name[:pos], append([]rune("**"), name[pos:]...)...)

				results = append(results, searchResult{name: string(name), id: elem.ID})
			}
		}
		db.RUnlock()
	}

	if len(results) == 0 {
		rsp.Message(db.Config.LangProperty("NoResults", nil))
		return
	}

	eodsort.Sort(results, len(results), func(index int) int {
		return results[index].id
	}, func(index int) string {
		return results[index].name
	}, func(index int, val string) {
		results[index].name = val
	}, sort, m.Author.ID, db, postfix)

	txt := make([]string, len(results))
	for i, val := range results {
		txt[i] = val.name
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      db.Config.LangProperty("ElemSearch", len(txt)),
		PageGetter: b.base.InvPageGetter,
		Items:      txt,
		User:       m.Author.ID,
	}, m, rsp)
}

func (b *Elements) Autocomplete(m types.Msg, query string) ([]string, types.GetResponse) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return nil, res
	}

	type searchResult struct {
		priority int
		id       int
	}
	results := make([]searchResult, 0)
	db.RLock()
	for _, el := range db.Elements {
		if strings.EqualFold(el.Name, query) {
			results = append(results, searchResult{0, el.ID})
		} else if strings.HasPrefix(strings.ToLower(el.Name), query) {
			results = append(results, searchResult{1, el.ID})
		} else if strings.Contains(strings.ToLower(el.Name), query) {
			results = append(results, searchResult{2, el.ID})
		}
		if len(results) > 1000 {
			break
		}
	}
	db.RUnlock()

	// sort by id
	sort.Slice(results, func(i, j int) bool {
		return results[i].id < results[j].id
	})
	// sort by priority
	sort.Slice(results, func(i, j int) bool {
		return results[i].priority < results[j].priority
	})
	if len(results) > types.AutocompleteResults {
		results = results[:types.AutocompleteResults]
	}

	names := make([]string, len(results))
	db.RLock()
	for i, res := range results {
		el, res := db.GetElement(res.id, true)
		if !res.Exists {
			continue
		}
		names[i] = el.Name
	}
	db.RUnlock()

	// sort by name
	sort.Strings(names)

	return names, types.GetResponse{Exists: true}
}
