package categories

import (
	"regexp"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

// flags: 0 = both, 1 = only categories, 2 = only vcats
func (b *Categories) Autocomplete(m types.Msg, query string, flags ...bool) ([]string, types.GetResponse) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return nil, res
	}

	type searchResult struct {
		priority int
		name     string
		size     int
	}
	results := make([]searchResult, 0)
	db.RLock()
	if len(flags) == 0 || len(flags) == 1 {
		for _, cat := range db.Cats() {
			if strings.EqualFold(cat.Name, query) {
				results = append(results, searchResult{0, cat.Name, len(cat.Elements)})
			} else if strings.HasPrefix(strings.ToLower(cat.Name), query) {
				results = append(results, searchResult{1, cat.Name, len(cat.Elements)})
			} else if strings.Contains(strings.ToLower(cat.Name), query) {
				results = append(results, searchResult{2, cat.Name, len(cat.Elements)})
			}
			if len(results) > 1000 {
				break
			}
		}
	}
	if len(flags) == 0 || len(flags) == 2 {
		for _, cat := range db.VCats() {
			if strings.EqualFold(cat.Name, query) {
				db.RUnlock()
				els, res := b.base.CalcVCat(cat, db, true)
				db.RLock()
				if res.Exists {
					results = append(results, searchResult{0, cat.Name, len(els)})
				}
			} else if strings.HasPrefix(strings.ToLower(cat.Name), query) {
				db.RUnlock()
				els, res := b.base.CalcVCat(cat, db, true)
				db.RLock()
				if res.Exists {
					results = append(results, searchResult{1, cat.Name, len(els)})
				}
			} else if strings.Contains(strings.ToLower(cat.Name), query) {
				db.RUnlock()
				els, res := b.base.CalcVCat(cat, db, true)
				db.RLock()
				if res.Exists {
					results = append(results, searchResult{2, cat.Name, len(els)})
				}
			}
			if len(results) > 1000 {
				break
			}
		}
	}
	db.RUnlock()

	// sort by length
	sort.Slice(results, func(i, j int) bool {
		return results[i].size > results[j].size
	})
	// sort by priority
	sort.Slice(results, func(i, j int) bool {
		return results[i].priority < results[j].priority
	})
	// shorten to max results
	if len(results) > types.AutocompleteResults {
		results = results[:types.AutocompleteResults]
	}

	names := make([]string, len(results))
	db.RLock()
	for i, res := range results {
		names[i] = res.name
	}
	db.RUnlock()

	// sort by name
	sort.Strings(names)

	return names, types.GetResponse{Exists: true}
}

func (b *Categories) SearchCmd(search string, sortKind string, regex bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	type searchResult struct {
		text  string
		name  string
		count int
	}

	// Make results
	results := make([]searchResult, 0)
	if regex {
		reg, err := regexp.Compile(search)
		if rsp.Error(err) {
			return
		}
		db.RLock()
		for _, cat := range db.Cats() {
			m := reg.FindIndex([]byte(cat.Name))
			if m != nil {
				name := []byte(cat.Name)
				name = append(name[:m[1]], append([]byte("**"), name[m[1]:]...)...)
				name = append(name[:m[0]], append([]byte("**"), name[m[0]:]...)...)

				results = append(results, searchResult{text: string(name), name: cat.Name, count: len(cat.Elements)})
			}
		}
		for _, cat := range db.VCats() {
			m := reg.FindIndex([]byte(cat.Name))
			if m != nil {
				name := []byte(cat.Name)
				name = append(name[:m[1]], append([]byte("**"), name[m[1]:]...)...)
				name = append(name[:m[0]], append([]byte("**"), name[m[0]:]...)...)

				count := 0
				if sortKind == "count" {
					v, res := b.base.CalcVCat(cat, db, true)
					if res.Exists {
						count = len(v)
					}
				}
				results = append(results, searchResult{text: string(name), name: cat.Name, count: count})
			}
		}
		db.RUnlock()
	} else {
		s := strings.ToLower(search)
		db.RLock()
		for _, cat := range db.Cats() {
			if strings.Contains(strings.ToLower(cat.Name), s) {
				name := []byte(cat.Name)
				pos := strings.Index(strings.ToLower(cat.Name), s)
				name = append(name[:pos+len(s)], append([]byte("**"), name[pos+len(s):]...)...)
				name = append(name[:pos], append([]byte("**"), name[pos:]...)...)

				results = append(results, searchResult{text: string(name), name: cat.Name, count: len(cat.Elements)})
			}
		}
		for _, cat := range db.VCats() {
			if strings.Contains(strings.ToLower(cat.Name), s) {
				name := []byte(cat.Name)
				pos := strings.Index(strings.ToLower(cat.Name), s)
				name = append(name[:pos+len(s)], append([]byte("**"), name[pos+len(s):]...)...)
				name = append(name[:pos], append([]byte("**"), name[pos:]...)...)

				count := 0
				if sortKind == "count" {
					v, res := b.base.CalcVCat(cat, db, true)
					if res.Exists {
						count = len(v)
					}
				}
				results = append(results, searchResult{text: string(name), name: cat.Name, count: count})
			}
		}
		db.RUnlock()
	}

	// Sort
	if len(results) == 0 {
		rsp.Message(db.Config.LangProperty("NoResults", nil))
		return
	}

	switch sortKind {
	case "count":
		sort.Slice(results, func(i, j int) bool { return results[i].count > results[j].count })

	default:
		sort.Slice(results, func(i, j int) bool { return eodsort.CompareStrings(results[i].name, results[j].name) })
	}

	txt := make([]string, len(results))
	for i, val := range results {
		txt[i] = val.text
	}
	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      "Category Search", // TODO: Translate
		PageGetter: b.base.InvPageGetter,
		Items:      txt,
		User:       m.Author.ID,
	}, m, rsp)
}
