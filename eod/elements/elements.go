package elements

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Elements struct {
	db    *sqlx.DB
	base  *base.Base
	polls *polls.Polls
	s     *sevcord.Sevcord
}

func (e *Elements) Init() {
	e.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"suggest",
		"Create a suggestion!",
		e.Suggest,
		sevcord.NewOption("result", "The result of the combination!", sevcord.OptionKindString, true),
		sevcord.NewOption("autocapitalize", "Whether or not to autocapitalize!", sevcord.OptionKindBool, false),
	))
	e.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"products",
		"View the elements that can be created using this element!",
		e.Products,
		sevcord.NewOption("element", "The element to view the products of!", sevcord.OptionKindInt, true).
			AutoComplete(e.Autocomplete),
	))
	e.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("edit", "Edit element properties!",
		sevcord.NewSlashCommand("name",
			"Edit the name of an element!",
			e.EditElementNameCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("name", "The new name of the element!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand("image",
			"Edit the image of an element!",
			e.EditElementImageCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("image", "The new image URL of the element!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand("description",
			"Edit the description of an element!",
			e.EditElementCommentCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("description", "The new description of the element!", sevcord.OptionKindString, true),
		),
		sevcord.NewSlashCommand("color",
			"Edit the color of an element!",
			e.EditElementColorCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("color", "The new color of the element! (decimal version of hex code)", sevcord.OptionKindInt, true),
		),
		sevcord.NewSlashCommand("creator",
			"Edit the creator of an element!",
			e.EditElementCreatorCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("creator", "The new creator of the element!", sevcord.OptionKindUser, true),
		),
		sevcord.NewSlashCommand("createdon",
			"Edit the creation date of an element!",
			e.EditElementCreatedonCmd,
			sevcord.NewOption("element", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
			sevcord.NewOption("createdon", "The new creation date of the element! (unix timestamp)", sevcord.OptionKindInt, true),
		),
	).
		RequirePermissions(discordgo.PermissionManageServer))
	e.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("delete", "Delete element properties!",
		sevcord.NewSlashCommand("combos",
			"Delete all combos apart from the first combo!",
			e.DeleteComboCmd,
			sevcord.NewOption("name", "The element to edit!", sevcord.OptionKindInt, true).
				AutoComplete(e.Autocomplete),
		),
	).
		RequirePermissions(discordgo.PermissionManageServer))
}

func NewElements(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base, polls *polls.Polls) *Elements {
	e := &Elements{
		db:    db,
		base:  base,
		polls: polls,
		s:     s,
	}
	e.Init()
	return e
}
