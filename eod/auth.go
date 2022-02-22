package eod

import (
	"github.com/Nv7-Github/Nv7Haven/eod/translation"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var noModCmds = map[string]types.Empty{
	"suggest":                  {},
	"mark":                     {},
	"image":                    {},
	"inv":                      {},
	"lb":                       {},
	"lbimage":                  {},
	"addcat":                   {},
	"cat":                      {},
	"hint":                     {},
	"stats":                    {},
	"idea":                     {},
	"help":                     {},
	"rmcat":                    {},
	"download":                 {},
	"elemsort":                 {},
	"breakdown":                {},
	"get":                      {},
	"setcolor":                 {},
	"invhint":                  {},
	"search":                   {},
	"View Inventory":           {},
	"View Info":                {},
	"Get Hint":                 {},
	"Get Inverse Hint":         {},
	"Get Color":                {},
	"View Leaderboard":         {},
	"View Inventory Breakdown": {},
	"color":                    {},
	"ai_idea":                  {},
	"wordcloud":                {},
	"delcat":                   {},
	"catop":                    {},
	"vcat":                     {},
}

func (b *EoD) canRunCmd(cmd *discordgo.InteractionCreate) (bool, string) {
	resp := cmd.ApplicationCommandData()

	// Check if mod is not required
	_, exists := noModCmds[resp.Name]
	if exists {
		return true, ""
	}

	// Check if is mod
	ismod, err := b.isMod(cmd.Member.User.ID, cmd.GuildID, b.newMsgSlash(cmd))
	if err != nil {
		return false, err.Error()
	}
	if ismod {
		return true, ""
	}

	// Get dat because everything after will require it
	db, res := b.GetDB(cmd.GuildID)
	if !res.Exists {
		return false, translation.LangProperty(translation.DefaultLang, "MustHaveAdmin", nil)
	}
	falseMsg := db.Config.LangProperty("MustHaveAdminOrModRole", db.Config.ModRole)

	// If command is path or catpath, check if has element/all elements in cat
	// path
	if resp.Name == "path" || resp.Name == "graph" || resp.Name == "notation" {
		if resp.Options[0].Name == "element" {
			inv := db.GetInv(cmd.Member.User.ID)
			if !res.Exists {
				return false, res.Message
			}

			el, res := db.GetElementByName(resp.Options[0].Options[0].StringValue())
			if !res.Exists {
				return true, "" // If the element doesn't exist, the cat command will tell the user it doesn't exist
			}

			exists = inv.Contains(el.ID)
			if !exists {
				return false, db.Config.LangProperty("MustHaveElemForPath", el.Name)
			}
			return true, ""
		} else {
			inv := db.GetInv(cmd.Member.User.ID)
			cat, res := db.GetCat(resp.Options[0].Options[0].StringValue())
			if !res.Exists {
				return true, "" // If the category doesn't exist, the cat command will tell the user it doesn't exist
			}

			// Check if user has all elements in category
			for elem := range cat.Elements {
				exists = inv.Contains(elem)
				if !exists {
					return false, db.Config.LangProperty("MustHaveCatForPath", cat.Name)
				}
			}

			return true, ""
		}
	}

	return false, falseMsg
}
