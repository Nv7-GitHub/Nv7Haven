package pages

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/queries"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
)

type Pages struct {
	base       *base.Base
	db         *sqlx.DB
	categories *categories.Categories
	elements   *elements.Elements
	queries    *queries.Queries
	s          *sevcord.Sevcord
}

func (p *Pages) PrintPostfix(postfixType string, elemName string, postfix string) (val string) {
	switch postfixType {
	case "length":
		return fmt.Sprintf(" - %d", len(elemName))
	case "creator", "colorer", "imager":
		return fmt.Sprintf(" - <@%s>", postfix)
	case "createdon":
		t, _ := time.Parse(time.RFC3339, postfix)
		return fmt.Sprintf(" - <t:%d>", t.Unix())
	case "treesize":
		val, _ := strconv.Atoi(postfix)
		return fmt.Sprintf(" - %s", humanize.Comma(int64(val)))
	case "color":
		val, _ := strconv.Atoi(postfix)
		return fmt.Sprintf(" - %s", util.FormatHex(val))
	default:
		return fmt.Sprintf(" - %s", postfix)
	}

}
func (p *Pages) Init() {
	// Inv
	p.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"inv",
		"View your inventory!",
		p.Inv,
		sevcord.NewOption("user", "The user to view the inventory of!", sevcord.OptionKindUser, false),
		sevcord.NewOption("sort", "The sort order of the inventory!", sevcord.OptionKindString, false).
			AddChoices(types.Sorts...),
		sevcord.NewOption("postfix", "Whether to add postfix!", sevcord.OptionKindBool, false),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	))
	p.s.AddButtonHandler("inv", p.InvHandler)

	// Lb
	p.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"lb",
		"View the leaderboard!",
		p.Lb,
		sevcord.NewOption("sort", "The sort order of the leaderboard!", sevcord.OptionKindString, false).
			AddChoices(lbSorts...),
		sevcord.NewOption("user", "The user to view the leaderboard from the point of view of!", sevcord.OptionKindUser, false),
		sevcord.NewOption("query", "View the stats within a query!", sevcord.OptionKindString, false).
			AutoComplete(p.queries.Autocomplete),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	))
	p.s.AddButtonHandler("lb", p.LbHandler)

	// Categories
	p.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("cat", "View categories!", sevcord.NewSlashCommand(
		"list",
		"View a list of all categories!",
		p.CatList,
		sevcord.NewOption("sort", "How to order the categories!", sevcord.OptionKindString, false).AddChoices(catListSorts...),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	), sevcord.NewSlashCommand(
		"view",
		"View a category's elements",
		p.Cat,
		sevcord.NewOption("category", "The category to view!", sevcord.OptionKindString, true).AutoComplete(p.categories.Autocomplete),
		sevcord.NewOption("sort", "How to order the elements!", sevcord.OptionKindString, false).AddChoices(types.Sorts...),
		sevcord.NewOption("postfix", "Whether to add postfix!", sevcord.OptionKindBool, false),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	), sevcord.NewSlashCommand(
		"add",
		"Add an element to a category!",
		p.categories.AddCat,
		sevcord.NewOption("category", "The category to add the elements to!", sevcord.OptionKindString, true).AutoComplete(p.categories.Autocomplete),
		sevcord.NewOption("element", "The element to add to the category!", sevcord.OptionKindInt, true).AutoComplete(p.elements.Autocomplete),
	), sevcord.NewSlashCommand(
		"remove",
		"Remove an element from a category!",
		p.categories.RmCat,
		sevcord.NewOption("category", "The category to remove elements from!", sevcord.OptionKindString, true).AutoComplete(p.categories.Autocomplete),
		sevcord.NewOption("element", "The element to remove from the category!", sevcord.OptionKindInt, true).AutoComplete(p.elements.Autocomplete),
	), sevcord.NewSlashCommand(
		"delete",
		"Delete all the elements from a category!",
		p.categories.DelCat,
		sevcord.NewOption("category", "The category to delete!", sevcord.OptionKindString, true).AutoComplete(p.categories.Autocomplete),
	), sevcord.NewSlashCommand(
		"addquery",
		"Add elements in a query to a category!",
		p.categories.AddCatQuery,
		sevcord.NewOption("category", "The category to add the elements to!", sevcord.OptionKindString, true).AutoComplete(p.categories.Autocomplete),
		sevcord.NewOption("query", "The query to add to the category!", sevcord.OptionKindString, true).AutoComplete(p.queries.Autocomplete),
	),
	))
	p.s.AddButtonHandler("catlist", p.CatListHandler)
	p.s.AddButtonHandler("cat", p.CatHandler)

	// Command lb
	p.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"commandlb",
		"See which commands are used the most!",
		p.CommandLb,
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	))
	p.s.AddButtonHandler("cmdlb", p.CommandLbHandler)

	// Queries
	p.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("query", "View queries!", sevcord.NewSlashCommand(
		"list",
		"View a list of every query!",
		p.QueryList,
		sevcord.NewOption("sort", "How to order the queries!", sevcord.OptionKindString, false).AddChoices(queryListSorts...),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	), sevcord.NewSlashCommand(
		"view",
		"View the elements in a query!",
		p.Query,
		sevcord.NewOption("query", "The query to view!", sevcord.OptionKindString, true).AutoComplete(p.queries.Autocomplete),
		sevcord.NewOption("sort", "How to sort the elements!", sevcord.OptionKindString, false).AddChoices(types.Sorts...),
		sevcord.NewOption("postfix", "Whether to add postfix!", sevcord.OptionKindBool, false),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	), sevcord.NewSlashCommand(
		"delete",
		"Delete a query!",
		p.queries.DeleteQuery,
		sevcord.NewOption("query", "The query to delete!", sevcord.OptionKindString, true).AutoComplete(p.queries.Autocomplete),
	), sevcord.NewSlashCommand(
		"download",
		"Download a query's contents!",
		p.queries.Download,
		sevcord.NewOption("query", "The query to download", sevcord.OptionKindString, true).
			AutoComplete(p.queries.Autocomplete),
		sevcord.NewOption("sort", "How to sort the elements!", sevcord.OptionKindString, false).
			AddChoices(types.Sorts...),
		sevcord.NewOption("postfix", "Retrieve data from the element!", sevcord.OptionKindString, false).
			AddChoices(types.Postfixes...),
	)))
	p.s.AddButtonHandler("querylist", p.QueryListHandler)
	p.s.AddButtonHandler("query", p.QueryHandler)

	// Products
	p.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"products",
		"View the elements that can be created using this element!",
		p.Products,
		sevcord.NewOption("element", "The element to view the products of!", sevcord.OptionKindInt, true).
			AutoComplete(p.elements.Autocomplete),
		sevcord.NewOption("sort", "How to order the elements!", sevcord.OptionKindString, false).AddChoices(types.Sorts...),
		sevcord.NewOption("postfix", "Whether to add postfix!", sevcord.OptionKindBool, false),
		sevcord.NewOption("page", "Which page of the results to view!", sevcord.OptionKindInt, false),
	))
	p.s.AddButtonHandler("products", p.ProductsHandler)
}

func NewPages(base *base.Base, db *sqlx.DB, s *sevcord.Sevcord, categories *categories.Categories, elements *elements.Elements, queries *queries.Queries) *Pages {
	p := &Pages{
		base:       base,
		db:         db,
		categories: categories,
		elements:   elements,
		queries:    queries,
		s:          s,
	}
	p.Init()
	return p
}
