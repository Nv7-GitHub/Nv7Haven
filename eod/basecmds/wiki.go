package basecmds

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func req(rsp types.Rsp, val any, url *url.URL) bool {
	res, err := http.Get(url.String())
	if rsp.Error(err) {
		return true
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(val)
	return rsp.Error(err)
}

type wikiSearchResultsVal struct {
	Query wikiSearchResults `json:"query"`
}

type wikiSearchResults struct {
	SearchInfo struct {
		TotalHits int `json:"totalhits"`
	} `json:"searchinfo"`
	Search []wikiSearchResult `json:"search"`
}

type wikiSearchResult struct {
	PageID int `json:"pageid"`
}

type wikiPages struct {
	Query struct {
		Pages map[int]wikiPage `json:"pages"`
	} `json:"query"`
}

type wikiPage struct {
	Extract string `json:"extract"`
}

// 1. https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=Nelson%20Mandela&format=json
// 2. https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exsentences=2&titles=Nelson%20Mandela&explaintext=1&format=json

func (b *BaseCmds) WikiCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get search results
	u := &url.URL{
		Scheme: "https",
		Host:   "en.wikipedia.org",
		Path:   "w/api.php",
	}
	q := u.Query()
	q.Set("action", "query")
	q.Set("list", "search")
	q.Set("srsearch", el.Name)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var results wikiSearchResultsVal
	if req(rsp, &results, u) {
		return
	}

	if results.Query.SearchInfo.TotalHits == 0 {
		rsp.ErrorMessage("No results found!") // TODO: Translate
		return
	}
	pageID := results.Query.Search[0].PageID

	// Get page
	u = &url.URL{
		Scheme: "https",
		Host:   "en.wikipedia.org",
		Path:   "w/api.php",
	}
	q = u.Query()
	q.Set("action", "query")
	q.Set("prop", "extracts")
	q.Set("exsentences", "2")
	q.Set("titles", el.Name)
	q.Set("explaintext", "1")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var pages wikiPages
	if req(rsp, &pages, u) {
		return
	}

	page := pages.Query.Pages[pageID]
	if len(strings.TrimSpace(page.Extract)) == 0 {
		rsp.ErrorMessage("No results found!") // TODO: Translate
		return
	}
	rsp.Message(page.Extract)
}
