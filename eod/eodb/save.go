package eodb

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (d *DB) SaveElement(el types.Element, new ...bool) error {
	d.Lock()
	defer d.Unlock()

	// Save to cache
	if len(new) > 0 {
		el.ID = len(d.Elements) + 1
		d.Elements = append(d.Elements, el)
		d.elemNames[strings.ToLower(el.Name)] = el.ID
	} else {
		d.Elements[el.ID-1] = el
	}

	// Persist
	dat, err := json.Marshal(el)
	if err != nil {
		return err
	}
	_, err = d.elemFile.WriteString(string(dat) + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) AddCombo(elems []int, result int) error {
	body := util.FormatCombo(elems)

	d.Lock()
	defer d.Unlock()

	d.combos[body] = result
	_, err := d.comboFile.WriteString(body + "=" + strconv.Itoa(result) + "\n")
	return err
}

func (d *DB) SaveConfig() error {
	d.Lock()
	defer d.Unlock()

	dat, err := json.Marshal(d.Config)
	if err != nil {
		return err
	}

	err = d.configFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = d.configFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = d.configFile.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) NewCat(name string) *types.Category {
	cat := &types.Category{
		Lock: &sync.RWMutex{},

		Name:     name,
		Guild:    d.Guild,
		Elements: make(map[int]types.Empty),
		Image:    "",
		Color:    0,
	}

	d.Lock()
	d.cats[strings.ToLower(name)] = cat
	d.Unlock()

	return cat
}

func (d *DB) SaveCat(elems *types.Category) error {
	// Empty?
	if len(elems.Elements) == 0 {
		d.Lock()
		delete(d.cats, strings.ToLower(elems.Name))
		delete(d.catFiles, strings.ToLower(elems.Name))
		d.Unlock()

		err := os.Remove(filepath.Join(d.dbPath, "categories", url.PathEscape(elems.Name)+".json"))
		return err
	}

	elems.Lock.RLock()
	dat, err := json.Marshal(elems)
	elems.Lock.RUnlock()
	if err != nil {
		return err
	}

	d.Lock()
	defer d.Unlock()

	file, exists := d.catFiles[strings.ToLower(elems.Name)]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "categories", url.PathEscape(elems.Name)+".json"))
		if err != nil {
			return err
		}
		d.catFiles[strings.ToLower(elems.Name)] = file
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SaveInv(inv *types.Inventory, recalc ...bool) error {
	inv.Lock.RLock()
	dat, err := json.Marshal(inv)
	inv.Lock.RUnlock()
	if err != nil {
		return err
	}

	d.RLock()
	if len(recalc) > 0 {
		for elem := range inv.Elements {
			elem, res := d.GetElement(elem, true)
			if !res.Exists {
				continue
			}
			if elem.Creator == inv.User {
				inv.MadeCnt++
			}
		}
	}
	d.RUnlock()

	file, exists := d.invFiles[inv.User]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "inventories", inv.User+".json"))
		if err != nil {
			return err
		}
		d.invFiles[inv.User] = file
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SavePoll(poll types.Poll) {
	d.Lock()
	d.Polls[poll.Message] = poll
	d.Unlock()
}

func (d *DB) DeletePoll(poll types.Poll) error {
	err := os.Remove(filepath.Join(d.dbPath, "polls", poll.Message+".json"))
	if err != nil {
		return err
	}

	d.Lock()
	delete(d.Polls, poll.Message)
	d.Unlock()

	return nil
}

func (d *DB) NewPoll(poll types.Poll) error {
	f, err := os.Create(filepath.Join(d.dbPath, "polls", poll.Message+".json"))
	if err != nil {
		return err
	}
	defer f.Close()

	dat, err := json.Marshal(poll)
	if err != nil {
		return err
	}
	_, err = f.Write(dat)
	if err != nil {
		return err
	}

	d.Lock()
	d.Polls[poll.Message] = poll
	d.Unlock()

	return nil
}
