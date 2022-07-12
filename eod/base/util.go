package base

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Base) GetColor(guild, id string) (int, error) {
	db, res := b.GetDB(guild)
	if res.Exists {
		db.Config.RLock()
		col, exists := db.Config.UserColors[id]
		db.Config.RUnlock()
		if exists {
			return col, nil
		}
	}

	mem, err := b.dg.State.Member(guild, id)
	if err != nil {
		mem, err = b.dg.GuildMember(guild, id)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	roles := make([]*discordgo.Role, len(mem.Roles))
	for i, roleID := range mem.Roles {
		role, err := b.GetRole(roleID, guild)
		if err != nil {
			return 0, err
		}
		roles[i] = role
	}

	sorted := discordgo.Roles(roles)
	sort.Sort(sorted)
	for _, role := range sorted {
		if role.Color != 0 {
			return role.Color, nil
		}
	}

	return 0, errors.New("eod: color not found")
}

func (b *Base) GetRole(id string, guild string) (*discordgo.Role, error) {
	role, err := b.dg.State.Role(guild, id)
	if err == nil {
		return role, nil
	}

	roles, err := b.dg.GuildRoles(guild)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.ID == id {
			return role, nil
		}
	}

	return nil, errors.New("eod: role not found")
}

type catSortInfo struct {
	Name string
	Cnt  int
}

func (b *Base) ElemCategories(elem int, db *eodb.DB, vcats bool) []string {
	// Get Categories
	cats := make([]catSortInfo, 0)
	db.RLock()
	for _, cat := range db.Cats() {
		_, exists := cat.Elements[elem]
		if exists {
			cats = append(cats, catSortInfo{
				Name: cat.Name,
				Cnt:  len(cat.Elements),
			})
		}
	}
	if vcats {
		for _, vcat := range db.VCats() {
			if fast && vcat.Rule == types.VirtualCategoryRuleSetOperation { // ignore set operations because they are slow
				continue
			}
			db.RUnlock()
			els, res := b.CalcVCat(vcat, db, true)
			db.RLock()
			if res.Exists {
				_, exists := els[elem]
				if exists {
					cats = append(cats, catSortInfo{
						Name: vcat.Name,
						Cnt:  len(els),
					})
				}
			}
		}
	}
	db.RUnlock()

	// Sort categories by count
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Cnt > cats[j].Cnt
	})

	// Convert to array
	out := make([]string, len(cats))
	for i, cat := range cats {
		out[i] = cat.Name
	}
	return out
}
