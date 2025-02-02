package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

var queryListSorts = []sevcord.Choice{
	sevcord.NewChoice("Name", "name"),
}

var queryListSortSql = map[string]string{
	"name": "name",
}

// Format: prevnext|sort|page
func (p *Pages) QueryListHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM queries WHERE guild=$1`, c.Guild()).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply pages
	page, _ := strconv.Atoi(parts[2])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var cats []struct {
		Name string `db:"name"`
	}
	err = p.db.Select(&cats, `SELECT name FROM queries WHERE guild=$1 ORDER BY `+queryListSortSql[parts[1]]+` LIMIT $2 OFFSET $3`, c.Guild(), length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for _, v := range cats {
		desc.WriteString(v.Name + "\n")
	}

	// Respond
	emb := sevcord.NewEmbed().
		Title(fmt.Sprintf("All Queries (%d)", cnt)).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple
	c.Respond(sevcord.NewMessage("").AddEmbed(emb).AddComponentRow(PageSwitchBtns("querylist", fmt.Sprintf("%s|%d", parts[1], page))...))
}

func (p *Pages) QueryList(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Params
	sort := "name"
	if opts[0] != nil {
		sort = opts[0].(string)
	}

	// Respond
	p.QueryListHandler(c, "next|"+sort+"|-1")
}

var queryPageCache = make(map[string]map[string]*types.Query)
var queryPageCacheLock = &sync.RWMutex{}

// Params: prevnext|user|sort|postfix|page|query
func (p *Pages) QueryHandler(c sevcord.Ctx, params string) {
	parts := strings.SplitN(params, "|", 6)

	// Get query
	var query *types.Query
	queryPageCacheLock.RLock()
	g, exists := queryPageCache[c.Guild()]
	if exists {
		query, exists = g[parts[5]]
	}
	queryPageCacheLock.RUnlock()

	if !exists {
		var ok bool
		query, ok = p.base.CalcQuery(c, parts[5])
		if !ok {
			return
		}
		queryPageCacheLock.Lock()
		v, exists := queryPageCache[c.Guild()]
		if !exists {
			v = make(map[string]*types.Query)
			queryPageCache[c.Guild()] = v
		}
		v[parts[5]] = query
		queryPageCacheLock.Unlock()
	}

	// Get count
	cnt := len(query.Elements)
	var common int
	err := p.db.QueryRow(`SELECT COALESCE(array_length($2 & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0)`, c.Guild(), pq.Array(query.Elements), parts[1]).Scan(&common)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[4])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var items []struct {
		Name    string `db:"name"`
		Cont    bool   `db:"cont"`
		Postfix string `db:"postfix"`
	}

	postfix := false
	if parts[3] == "1" {
		postfix = true
	} else {
		postfix = false
	}
	//false if not valid in DB
	postfixable := parts[2] != "length" && parts[2] != "found"
	if postfix && postfixable {
		err = p.db.Select(&items, `SELECT name, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont, `+parts[2]+` postfix FROM elements WHERE id=ANY($2) AND guild=$1 ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), pq.Array(query.Elements), length, length*page, parts[1])

	} else {
		err = p.db.Select(&items, `SELECT name, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE id=ANY($2) AND guild=$1 ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), pq.Array(query.Elements), length, length*page, parts[1])

	}

	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for _, v := range items {
		if v.Cont {
			fmt.Fprintf(desc, "%s %s", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s %s", v.Name, types.NoCheck)
		}
		if postfix && parts[2] != "found" {
			desc.WriteString(p.PrintPostfix(parts[2], v.Name, v.Postfix))

		}
		desc.WriteString("\n")

	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s (%s, %s%%)", parts[5], humanize.Comma(int64(cnt)), humanize.FormatFloat("", float64(common)/float64(cnt)*100))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("query", fmt.Sprintf("%s|%s|%s|%d|%s", parts[1], parts[2], parts[3], page, parts[5]))...),
	)
}

func (p *Pages) Query(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	// Get params
	sort := "id"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Get name
	var name string
	err := p.db.QueryRow(`SELECT name FROM queries WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(args[0].(string))).Scan(&name)
	if err != nil {
		p.base.Error(c, err, "Query **"+args[0].(string)+"** doesn't exist!")
		return
	}

	// Reset if its there
	queryPageCacheLock.Lock()
	g, exists := queryPageCache[c.Guild()]
	if exists {
		delete(g, name)
	}
	queryPageCacheLock.Unlock()
	postfix := false
	postfixval := 0
	if args[2] != nil {
		postfix = args[2].(bool)
	}
	if postfix {
		postfixval = 1
	} else {
		postfixval = 0
	}
	// Create embed
	p.QueryHandler(c, fmt.Sprintf("next|%s|%s|%d|-1|%s", c.Author().User.ID, sort, postfixval, name))
}
