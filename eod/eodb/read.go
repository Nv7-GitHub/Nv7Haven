package eodb

import (
	"fmt"
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
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf(d.Config.LangProperty("DoesntExist"), name),
		}
	}
	return d.Elements[id-1], types.GetResponse{Exists: true}
}

func (d *DB) GetIDByName(name string) (int, types.GetResponse) {
	d.RLock()
	defer d.RUnlock()

	id, exists := d.elemNames[strings.ToLower(name)]
	if !exists {
		return 0, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf(d.Config.LangProperty("DoesntExist"), name),
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
				Message: fmt.Sprintf(d.Config.LangProperty("DoesntExist"), "#0"),
			}
		}
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("IDCannotBeNegative"),
		}
	}
	if id > len(d.Elements) {
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf(d.Config.LangProperty("DoesntExist"), "#"+strconv.Itoa(id)),
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
			Message: d.Config.LangProperty("ComboNoExist"),
		}
	}
	return res, types.GetResponse{Exists: true}
}

func (d *DB) GetInv(id string) *types.Inventory {
	d.RLock()
	inv, exists := d.invs[id]
	d.RUnlock()
	if !exists {
		inv = types.NewInventory(id, map[int]types.Empty{
			1: {},
			2: {},
			3: {},
			4: {},
		}, 0)
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
			Message: fmt.Sprintf(d.Config.LangProperty("CatNoExist"), name),
		}
	}
	return cat, types.GetResponse{Exists: true}
}

func (d *DB) GetPoll(id string) (types.Poll, types.GetResponse) {
	d.RLock()
	poll, exists := d.Polls[id]
	d.RUnlock()
	if !exists {
		return types.Poll{}, types.GetResponse{
			Exists:  false,
			Message: d.Config.LangProperty("PollNoExist"),
		}
	}
	return poll, types.GetResponse{Exists: true}
}
