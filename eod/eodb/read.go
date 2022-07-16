package eodb

import (
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (d *DB) GetElementByName(name string, nolock ...bool) (types.Element, types.GetResponse) {
	if len(nolock) == 0 {
		d.RLock()
		defer d.RUnlock()
	}

	id, exists := d.elemNames[strings.ToLower(name)]
	if !exists {
		if len(name) > 1 && name[0] == '#' {
			id, err := strconv.Atoi(name[1:])
			if err == nil {
				return d.GetElement(id)
			}
		}

		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("DoesntExist", name),
		}
	}
	return d.Elements[id-1], types.GetResponse{Exists: true}
}

func (d *DB) GetIDByName(name string) (int, types.GetResponse) {
	d.RLock()
	defer d.RUnlock()

	id, exists := d.elemNames[strings.ToLower(name)]
	if !exists {
		if name[0] == '#' && len(name) > 1 {
			id, err := strconv.Atoi(name[1:])
			if err == nil {
				// Code from GetElement
				if id < 1 {
					if id == 0 {
						return 0, types.GetResponse{
							Exists:  false,
							Message: d.Config.LangProperty("DoesntExist", "#0"),
						}
					}
					return 0, types.GetResponse{
						Exists:  false,
						Message: d.Config.LangProperty("IDCannotBeNegative", nil),
					}
				}
				if id > len(d.Elements) {
					return 0, types.GetResponse{
						Exists:  false,
						Message: d.Config.LangProperty("DoesntExist", "#"+strconv.Itoa(id)),
					}
				}

				return id, types.GetResponse{Exists: true}
			}
		}

		return 0, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("DoesntExist", name),
		}
	}
	return id, types.GetResponse{Exists: true}
}

func (d *DB) GetElement(id int, nolock ...bool) (types.Element, types.GetResponse) {
	if len(nolock) == 0 {
		d.RLock()
		defer d.RUnlock()
	}

	if id < 1 {
		if id == 0 {
			return types.Element{}, types.GetResponse{
				Exists:  false,
				Message: d.Config.LangProperty("DoesntExist", "#0"),
			}
		}
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("IDCannotBeNegative", nil),
		}
	}
	if id > len(d.Elements) {
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("DoesntExist", "#"+strconv.Itoa(id)),
		}
	}

	return d.Elements[id-1], types.GetResponse{Exists: true}
}

func (d *DB) GetCombo(elems []int) (int, types.GetResponse) {
	txt := util.FormatCombo(elems)
	d.RLock()
	res, exists := d.combos[txt]
	d.RUnlock()
	if !exists {
		return 0, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("DBNoCombo", nil),
		}
	}
	return res, types.GetResponse{Exists: true}
}

func defaultInv(id string) *types.Inventory {
	return types.NewInventory(id, map[int]types.Empty{
		1: {},
		2: {},
		3: {},
		4: {},
	}, 0)
}

func (d *DB) GetInv(id string) *types.Inventory {
	d.RLock()
	inv, exists := d.invs[id]
	d.RUnlock()
	if !exists {
		inv = defaultInv(id)
		d.Lock()
		d.invs[id] = inv
		d.Unlock()

		return inv
	}
	return inv
}

func (d *DB) GetCat(name string) (*types.Category, types.GetResponse) {
	d.RLock()
	cat, exists := d.cats[strings.ToLower(name)]
	d.RUnlock()
	if !exists {
		return nil, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("CatNoExist", name),
		}
	}
	return cat, types.GetResponse{Exists: true}
}

func (d *DB) GetCatCache(name string) (map[int]types.Empty, bool) {
	d.RLock()
	cache, exists := d.catCache[strings.ToLower(name)]
	d.RUnlock()
	if !exists {
		return nil, false
	}
	return cache, true
}

func (d *DB) GetVCat(name string) (*types.VirtualCategory, types.GetResponse) {
	d.RLock()
	vcat, exists := d.vcats[strings.ToLower(name)]
	d.RUnlock()
	if !exists {
		return nil, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("CatNoExist", name),
		}
	}
	return vcat, types.GetResponse{Exists: true}
}

func (d *DB) GetPoll(id string) (types.Poll, types.GetResponse) {
	d.RLock()
	poll, exists := d.Polls[id]
	d.RUnlock()
	if !exists {
		return types.Poll{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("PollNoExist", nil),
		}
	}
	return poll, types.GetResponse{Exists: true}
}
