package eod

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var noModCmds = map[string]types.Empty{
	"suggest":          {},
	"mark":             {},
	"image":            {},
	"inv":              {},
	"lb":               {},
	"addcat":           {},
	"cat":              {},
	"hint":             {},
	"stats":            {},
	"idea":             {},
	"help":             {},
	"rmcat":            {},
	"download":         {},
	"elemsort":         {},
	"breakdown":        {},
	"get":              {},
	"setcolor":         {},
	"invhint":          {},
	"elemsearch":       {},
	"Get Inventory":    {},
	"Get Info":         {},
	"Get Hint":         {},
	"Get Inverse Hint": {},
	"Get Color":        {},
	"View Leaderboard": {},
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
	lock.RLock()
	dat, exists := b.dat[cmd.GuildID]
	lock.RUnlock()
	falseMsg := "You need to have permission `Administrator` or have role <@&" + dat.ModRole + ">!"
	if !exists {
		return false, falseMsg
	}

	// If command is path or catpath, check if has element/all elements in cat
	// path
	if resp.Name == "path" || resp.Name == "graph" || resp.Name == "notation" {
		if resp.Options[0].Name == "element" {
			inv, res := dat.GetInv(cmd.Member.User.ID, true)
			if !res.Exists {
				return false, res.Message
			}

			el, res := dat.GetElement(resp.Options[0].Options[0].StringValue())
			if !res.Exists {
				return true, "" // If the element doesn't exist, the cat command will tell the user it doesn't exist
			}

			exists = inv.Contains(el.Name)
			if !exists {
				return false, fmt.Sprintf("You must have element **%s** to get it's path!", el.Name)
			}
			return true, ""
		} else {
			inv, res := dat.GetInv(cmd.Member.User.ID, true)
			if !res.Exists {
				return false, res.Message
			}

			cat, res := dat.GetCategory(resp.Options[0].Options[0].StringValue())
			if !res.Exists {
				return true, "" // If the category doesn't exist, the cat command will tell the user it doesn't exist
			}

			// Check if user has all elements in category
			for elem := range cat.Elements {
				exists = inv.Contains(elem)
				if !exists {
					return false, fmt.Sprintf("You must have all elements in category **%s** to get its path!", cat.Name)
				}
			}

			return true, ""
		}
	}

	return false, falseMsg
}
