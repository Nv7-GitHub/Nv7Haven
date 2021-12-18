package elements

import (
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const resultCount = 25

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
	if len(results) > resultCount {
		results = results[:resultCount]
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
