package eod

import (
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/pages"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/queries"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Init() {
	b.base = base.NewBase(b.s, b.db)
	b.polls = polls.NewPolls(b.db, b.base, b.s)
	b.elements = elements.NewElements(b.s, b.db, b.base, b.polls)
	b.categories = categories.NewCategories(b.db, b.base, b.s, b.polls)
	b.queries = queries.NewQueries(b.s, b.db, b.base, b.polls, b.elements, b.categories)
	b.pages = pages.NewPages(b.base, b.db, b.s, b.categories, b.elements, b.queries)
	b.s.SetMessageHandler(b.messageHandler)

	// Start saving stats
	go func() {
		b.base.SaveStats()
		for {
			time.Sleep(time.Minute * 30)
			b.base.SaveStats()
		}
	}()

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
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("image", "Change an image!",
		sevcord.NewSlashCommand(
			"element",
			"Change the image of an element!",
			func(c sevcord.Ctx, opts []any) {
				b.elements.ImageCmd(c, int(opts[0].(int64)), opts[1].(*sevcord.SlashCommandAttachment).URL)
			},
			sevcord.NewOption("element", "The element to change the image of!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
			sevcord.NewOption("image", "The image to change it to!", sevcord.OptionKindAttachment, true),
		),
		sevcord.NewSlashCommand(
			"category",
			"Change the image of a category!",
			func(c sevcord.Ctx, opts []any) {
				b.categories.ImageCmd(c, opts[0].(string), opts[1].(*sevcord.SlashCommandAttachment).URL)
			},
			sevcord.NewOption("category", "The category to change the image of!", sevcord.OptionKindString, true).
				AutoComplete(b.categories.Autocomplete),
			sevcord.NewOption("image", "The image to change it to!", sevcord.OptionKindAttachment, true),
		),
		sevcord.NewSlashCommand(
			"query",
			"Change the image of a query!",
			func(c sevcord.Ctx, opts []any) {
				b.queries.ImageCmd(c, opts[0].(string), opts[1].(*sevcord.SlashCommandAttachment).URL)
			},
			sevcord.NewOption("query", "The query to change the image of!", sevcord.OptionKindString, true).
				AutoComplete(b.queries.Autocomplete),
			sevcord.NewOption("image", "The image to change it to!", sevcord.OptionKindAttachment, true),
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
			sevcord.NewOption("category", "The category to change the color of!", sevcord.OptionKindString, true).
				AutoComplete(b.categories.Autocomplete),
			sevcord.NewOption("color", "The hex code of the color to change it to!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand(
			"query",
			"Change the color of a query!",
			b.queries.ColorCmd,
			sevcord.NewOption("query", "The query to change the color of!", sevcord.OptionKindString, true).
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
		sevcord.NewSlashCommand(
			"categories",
			"See the categories an element is in!",
			b.pages.ElemCats,
			sevcord.NewOption("element", "The element to view the categories of!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
		),
		sevcord.NewSlashCommand(
			"found",
			"See who has found an element!",
			b.pages.ElemFound,
			sevcord.NewOption("element", "The element to view the people who have found!", sevcord.OptionKindInt, true).
				AutoComplete(b.elements.Autocomplete),
		),
	))
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("edit", "Edit element properties!",
		sevcord.NewSlashCommandGroup("element", "Edit element properties!",
			sevcord.NewSlashCommand("name",
				"Edit the name of an element!",
				b.elements.EditElementNameCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("name", "The new name of the element!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand("image",
				"Edit the image of an element!",
				b.elements.EditElementImageCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("image", "The new image URL of the element!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand("description",
				"Edit the description of an element!",
				b.elements.EditElementCommentCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("description", "The new description of the element!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand("color",
				"Edit the color of an element!",
				b.elements.EditElementColorCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("color", "The new color of the element! (decimal version of hex code)", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand("creator",
				"Edit the creator of an element!",
				b.elements.EditElementCreatorCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("creator", "The new creator of the element!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand("createdon",
				"Edit the creation date of an element!",
				b.elements.EditElementCreatedonCmd,
				sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
					AutoComplete(b.elements.Autocomplete),
				sevcord.NewOption("createdon", "The new creation date of the element! (unix timestamp)", sevcord.OptionKindInt, true),
			),
		),
		sevcord.NewSlashCommandGroup("query", "Edit the properties of elements in a query!",
			sevcord.NewSlashCommand("image",
				"Edit the images of elements in a query!",
				b.queries.EditElementImageCmd,
				sevcord.NewOption("query", "The query to edit!", sevcord.OptionKindString, true).
					AutoComplete(b.queries.Autocomplete),
				sevcord.NewOption("image", "The new image URL of the element!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand("description",
				"Edit the description of elements in a query!",
				b.queries.EditElementCommentCmd,
				sevcord.NewOption("query", "The query to edit!", sevcord.OptionKindString, true).
					AutoComplete(b.queries.Autocomplete),
				sevcord.NewOption("description", "The new description of the element!", sevcord.OptionKindString, true),
			),
			sevcord.NewSlashCommand("color",
				"Edit the color of elements in a query!",
				b.queries.EditElementColorCmd,
				sevcord.NewOption("query", "The query to edit!", sevcord.OptionKindString, true).
					AutoComplete(b.queries.Autocomplete),
				sevcord.NewOption("color", "The new color of the element! (decimal version of hex code)", sevcord.OptionKindInt, true),
			),
			sevcord.NewSlashCommand("creator",
				"Edit the creator of an elements in a query!",
				b.queries.EditElementCreatorCmd,
				sevcord.NewOption("query", "The query to edit!", sevcord.OptionKindString, true).
					AutoComplete(b.queries.Autocomplete),
				sevcord.NewOption("creator", "The new creator of the element!", sevcord.OptionKindUser, true),
			),
			sevcord.NewSlashCommand("createdon",
				"Edit the creation date of elements in a query!",
				b.queries.EditElementCreatedonCmd,
				sevcord.NewOption("query", "The query to edit!", sevcord.OptionKindString, true).
					AutoComplete(b.queries.Autocomplete),
				sevcord.NewOption("createdon", "The new creation date of the element! (unix timestamp)", sevcord.OptionKindInt, true),
			),
		),
	).
		RequirePermissions(discordgo.PermissionManageServer))
	b.s.AddButtonHandler("elemcats", b.pages.ElemCatHandler)
	b.s.AddButtonHandler("elemfound", b.pages.ElemFoundHandler)
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
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"next",
		"Find the next element to make!",
		b.elements.Next,
		sevcord.NewOption("query", "A query to select the random element to be made from!", sevcord.OptionKindString, false).
			AutoComplete(b.queries.Autocomplete),
	))
	b.s.AddButtonHandler("next", b.elements.NextHandler)
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"idea",
		"Get an element idea!",
		b.elements.Idea,
		sevcord.NewOption("query", "A query to select the elements in the idea to be made from!", sevcord.OptionKindString, false).
			AutoComplete(b.queries.Autocomplete),
		sevcord.NewOption("count", "The number of elements to include in the idea!", sevcord.OptionKindInt, false).
			MinMax(2, types.MaxComboLength),
	))
	b.s.AddButtonHandler("idea", b.elements.IdeaHandler)
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand("uncheese", "Delete the lowest tree size combo for elements!", b.elements.Uncheese, sevcord.NewOption("query", "A query to select the elements from!", sevcord.OptionKindString, true).
		AutoComplete(b.queries.Autocomplete)).
		RequirePermissions(discordgo.PermissionManageServer))
	b.s.AddButtonHandler("uncheese", b.elements.UncheeseHandler)
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"ping",
		"Check the server ping!",
		b.PingCmd,
	))
}
