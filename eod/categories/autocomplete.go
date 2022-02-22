package categories

import (
	"sort"
	"strings"

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
				els, res := b.base.CalcVCat(cat, db)
				if res.Exists {
					results = append(results, searchResult{0, cat.Name, len(els)})
				}
			} else if strings.HasPrefix(strings.ToLower(cat.Name), query) {
				els, res := b.base.CalcVCat(cat, db)
				if res.Exists {
					results = append(results, searchResult{1, cat.Name, len(els)})
				}
			} else if strings.Contains(strings.ToLower(cat.Name), query) {
				els, res := b.base.CalcVCat(cat, db)
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
