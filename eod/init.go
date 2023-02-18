package eod

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/pages"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/queries"
	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Bot) Init() {
	b.base = base.NewBase(b.s, b.db)
	b.polls = polls.NewPolls(b.db, b.base, b.s)
	b.elements = elements.NewElements(b.s, b.db, b.base, b.polls)
	b.categories = categories.NewCategories(b.db, b.base, b.s, b.polls)
	b.queries = queries.NewQueries(b.s, b.db, b.base, b.polls, b.elements, b.categories)
	b.pages = pages.NewPages(b.base, b.db, b.s, b.categories, b.elements, b.queries)
	b.s.SetMessageHandler(b.messageHandler)

	// Commands
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("sign", "Change a comment!",
		sevcord.NewSlashCommand(
			"element",
			"Change the comment of an element!",
			b.elements.SignCmd,
			sevcord.NewOption("element", "The element to change the comment of!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"category",
			"Change the comment of a category!",
			b.categories.SignCmd,
			sevcord.NewOption("category", "The category to change the comment of!", sevcord.OptionKindString, true).
				AutoComplete(b.categories.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"query",
			"Change the comment of a query!",
			b.queries.SignCmd,
			sevcord.NewOption("query", "The query to change the comment of!", sevcord.OptionKindString, true).
				AutoComplete(b.queries.Autocomplete),
		),
	))
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("color", "Change a color!",
		sevcord.NewSlashCommand(
			"element",
			"Change the color of an element!",
			b.elements.ColorCmd,
			sevcord.NewOption("element", "The element to change the color of!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
			sevcord.NewOption("color", "The hex code of the color to change it to!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand(
			"category",
			"Change the color of a category!",
			b.categories.ColorCmd,
			sevcord.NewOption("category", "The category to change the image of!", sevcord.OptionKindString, true).
				AutoComplete(b.categories.Autocomplete),
			sevcord.NewOption("color", "The hex code of the color to change it to!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand(
			"query",
			"Change the color of a query!",
			b.queries.ColorCmd,
			sevcord.NewOption("query", "The query to change the image of!", sevcord.OptionKindString, true).
				AutoComplete(b.queries.Autocomplete),
			sevcord.NewOption("color", "The hex code of the color to change it to!", sevcord.OptionKindString, true),
		),
	))
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("info", "Get element, category, or query info!",
		sevcord.NewSlashCommand(
			"element",
			"Get element info!",
			b.elements.InfoSlashCmd,
			sevcord.NewOption("element", "The element to view the info of!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"category",
			"Get category info!",
			b.categories.Info,
			sevcord.NewOption("category", "The category to view the info of!", sevcord.OptionKindString, true).
				AutoComplete(b.categories.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"query",
			"Get query info!",
			b.queries.Info,
			sevcord.NewOption("query", "The query to view the info of!", sevcord.OptionKindString, true).
				AutoComplete(b.queries.Autocomplete),
		),
	))
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"hint",
		"Learn how to make an element!",
		b.elements.Hint,
		sevcord.NewOption("element", "An element to get the hint of!", sevcord.OptionKindInt, false).
			AutoComplete(b.elements.Autocomplete),
		sevcord.NewOption("query", "A query to select the random element to be made from!", sevcord.OptionKindString, false).
			AutoComplete(b.queries.Autocomplete),
	))
	b.s.AddButtonHandler("hint", b.elements.HintHandler)
}
