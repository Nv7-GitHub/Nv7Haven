package eod

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var noModCmds = map[string]empty{
	"suggest":      {},
	"mark":         {},
	"image":        {},
	"inv":          {},
	"lb":           {},
	"addcat":       {},
	"cat":          {},
	"hint":         {},
	"stats":        {},
	"idea":         {},
	"about":        {},
	"rmcat":        {},
	"catimg":       {},
	"downloadinv":  {},
	"elemsort":     {},
	"breakdown":    {},
	"catbreakdown": {},
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
	falseMsg := "You need to have permission `Administrator` or have role <@&" + dat.modRole + ">!"
	if !exists {
		return false, falseMsg
	}

	// If command is path or catpath, check if has element/all elements in cat
	// path
	if resp.Name == "path" {
		dat.lock.RLock()
		inv, exists := dat.invCache[cmd.Member.User.ID]
		dat.lock.RUnlock()
		if !exists {
			return false, "You don't have an inventory!"
		}

		name := strings.ToLower(resp.Options[0].StringValue())
		dat.lock.RLock()
		el, exists := dat.elemCache[name]
		dat.lock.RUnlock()
		if !exists {
			return true, "" // If the element doesn't exist, the cat command will tell the user it doesn't exist
		}

		_, exists = inv[name]
		if !exists {
			return false, fmt.Sprintf("You must have element **%s** to get it's path!", el.Name)
		}
		return true, ""
	}

	// catpath
	if resp.Name == "catpath" {
		dat.lock.RLock()
		inv, exists := dat.invCache[cmd.Member.User.ID]
		dat.lock.RUnlock()
		if !exists {
			return false, "You don't have an inventory!"
		}
		cat, exists := dat.catCache[strings.ToLower(resp.Options[0].StringValue())]
		if !exists {
			return true, "" // If the category doesn't exist, the cat command will tell the user it doesn't exist
		}

		// Check if user has all elements in category
		for elem := range cat.Elements {
			_, exists = inv[strings.ToLower(elem)]
			if !exists {
				return false, fmt.Sprintf("You must have all elements in category **%s** to get its path!", cat.Name)
			}
		}

		return true, ""
	}

	return false, falseMsg
}
