package queries

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db         *sqlx.DB
	base       *base.Base
	s          *sevcord.Sevcord
	polls      *polls.Polls
	elements   *elements.Elements
	categories *categories.Categories
}

func (q *Queries) Init() {
	q.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("newquery", "Create a new query!",
		sevcord.NewSlashCommand(
			"element",
			"Create a query that contains a single element!",
			q.CreateElementCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("element", "The element for the query to contain!", sevcord.OptionKindInt, true).
				AutoComplete(q.elements.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"category",
			"Create a query that contains the elements in a category!",
			q.CreateCategoryCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("category", "The category for the query to contain!", sevcord.OptionKindString, true).
				AutoComplete(q.categories.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"products",
			"Create a query that contains the products of every element in another query!",
			q.CreateProductsCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("query", "The query to contain the products of!", sevcord.OptionKindString, true).
				AutoComplete(q.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"parents",
			"Create a query that contains the parents of every element in another query!",
			q.CreateParentsCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("query", "The query to contain the parents of!", sevcord.OptionKindString, true).
				AutoComplete(q.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"inventory",
			"Create a query that contains the elements in a user's inventory!",
			q.CreateInventoryCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("user", "The query to contain the inventory of!", sevcord.OptionKindUser, true),
		),
		sevcord.NewSlashCommand(
			"regex",
			"Create a query that filters the names of elements in a query by a POSIX-style regex!",
			q.CreateRegexCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("query", "The query to filter!", sevcord.OptionKindString, true).
				AutoComplete(q.Autocomplete),
			sevcord.NewOption("regex", "The regex to filter by!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand(
			"elements",
			"Create a query that contains every element!",
			q.CreateElementsCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommandGroup(
			"comparison",
			"Create a query that compares the elements in a query!",
			sevcord.NewSlashCommand(
				"id",
				"Compare the IDs of the elements!",
				q.CreateComparisonIDCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The ID to compare by!", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand(
				"name",
				"Compare the names of the elements!",
				q.CreateComparisonNameCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The name to compare by!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand(
				"image",
				"Compare the images of the elements!",
				q.CreateComparisonImageCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The image to compare by!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand(
				"color",
				"Compare the colors of the elements!",
				q.CreateComparisonColorCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The hex code of the color to compare by!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand(
				"description",
				"Compare the descriptions of the elements!",
				q.CreateComparisonDescriptionCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The description to compare by!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand(
				"creator",
				"Compare the creators of the elements!",
				q.CreateComparisonCreatorCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The creator to compare by!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand(
				"commenter",
				"Compare the commeneters of the elements!",
				q.CreateComparisonCommenterCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The commenter to compare by!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand(
				"colorer",
				"Compare the colorers of the elements!",
				q.CreateComparisonColorerCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The colorer to compare by!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand(
				"imager",
				"Compare the imagers of the elements!",
				q.CreateComparisonImagerCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The imager to compare by!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand(
				"treesize",
				"Compare the tree size of the elements!",
				q.CreateComparisonTreesizeCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The treesize to compare by!", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand(
				"length",
				"Compare the length of the names of the elements!",
				q.CreateComparisonLengthCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The length to compare by!", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand(
				"usedin",
				"Compare the used in of the elements!",
				q.CreateComparisonUsedinCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The used in to compare by!", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand(
				"madewith",
				"Compare the made with of the elements!",
				q.CreateComparisonMadeWithCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The made with to compare by!", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand(
				"tier",
				"Compare the tier of the elements!",
				q.CreateComparisonTierCmd,
				sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
				sevcord.NewOption("operator", "The operator to compare by!", sevcord.OptionKindString, true).AddChoices(ComparisonQueryOpChoices...),
				sevcord.NewOption("value", "The tier to compare by!", sevcord.OptionKindInt, true),
			),
		),
		sevcord.NewSlashCommand(
			"operation",
			"Create a query that performs a set operation on two queries!",
			q.CreateOperationCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
			sevcord.NewOption("left", "The left side of the operation!", sevcord.OptionKindString, true).
				AutoComplete(q.Autocomplete),
			sevcord.NewOption("right", "The right side of the operation!", sevcord.OptionKindString, true).
				AutoComplete(q.Autocomplete),
			sevcord.NewOption("operator", "The operation to perform!", sevcord.OptionKindString, true).
				AddChoices(sevcord.NewChoice("Union", "union"),
					sevcord.NewChoice("Intersection", "intersection"),
					sevcord.NewChoice("Difference", "difference"),
				),
		),
	))
	q.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"path",
		"Learn how to make the elements in a query!",
		func(c sevcord.Ctx, opts []any) {
			q.PathCmd(c, opts, true)
		},
		sevcord.NewOption("query", "The query to view the path of!", sevcord.OptionKindString, true).
			AutoComplete(q.Autocomplete),
	))
	q.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"pathjson",
		"Learn how to make the elements in a query!",
		func(c sevcord.Ctx, opts []any) {
			q.PathCmd(c, opts, false)
		},
		sevcord.NewOption("query", "The query to view the path of!", sevcord.OptionKindString, true).
			AutoComplete(q.Autocomplete),
	).RequirePermissions(discordgo.PermissionManageServer))
	q.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"give",
		"Give the elements in a query to a user!",
		q.base.Give,
		sevcord.NewOption("user", "The user to give the elements to!", sevcord.OptionKindUser, true),
		sevcord.NewOption("query", "The query to give the elements of!", sevcord.OptionKindString, true).AutoComplete(q.Autocomplete),
	).RequirePermissions(discordgo.PermissionManageChannels))
}

func NewQueries(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base, polls *polls.Polls, elements *elements.Elements, categories *categories.Categories) *Queries {
	q := &Queries{
		db:         db,
		base:       base,
		s:          s,
		polls:      polls,
		elements:   elements,
		categories: categories,
	}
	q.Init()
	return q
}
