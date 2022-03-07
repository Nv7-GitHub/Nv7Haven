package eodb

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/sasha-s/go-deadlock"
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
		old := d.Elements[el.ID-1]
		if old.Name != el.Name {
			delete(d.elemNames, strings.ToLower(old.Name))
			d.elemNames[strings.ToLower(el.Name)] = el.ID
		}
		d.Elements[el.ID-1] = el
	}

	if d.inTransaction { // Don't persist
		return nil
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
	// Add to AI
	d.AI.AddCombo(body, false)

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

	_, err = d.configFile.Seek(0, 0)
	if err != nil {
		return err
	}
	err = d.configFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = d.configFile.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SaveVCat(vcat *types.VirtualCategory) error {
	d.Lock()
	defer d.Unlock()

	d.vcats[strings.ToLower(vcat.Name)] = vcat

	dat, err := json.Marshal(d.vcats)
	if err != nil {
		return err
	}

	_, err = d.vcatsFile.Seek(0, 0)
	if err != nil {
		return err
	}
	err = d.vcatsFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = d.vcatsFile.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteVCat(name string) error {
	d.Lock()
	defer d.Unlock()

	delete(d.vcats, strings.ToLower(name))

	dat, err := json.Marshal(d.vcats)
	if err != nil {
		return err
	}

	_, err = d.vcatsFile.Seek(0, 0)
	if err != nil {
		return err
	}
	err = d.vcatsFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = d.vcatsFile.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) NewCat(name string) *types.Category {
	cat := &types.Category{
		Lock: &deadlock.RWMutex{},

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

func (d *DB) SaveCatCache(name string, elems map[int]types.Empty) error {
	d.Lock()
	defer d.Unlock()

	file, exists := d.catCacheFiles[strings.ToLower(name)]
	var err error
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "catcache", url.PathEscape(name)+".json"))
		if err != nil {
			return err
		}
		d.catCacheFiles[strings.ToLower(name)] = file
	}
	cache, exists := d.catCache[strings.ToLower(name)]
	if !exists {
		cache = make(map[int]types.Empty)
	}

	// Calc diff
	rm := make(map[int]types.Empty)
	add := make(map[int]types.Empty)
	for el := range elems {
		_, exists := cache[el]
		if !exists {
			add[el] = types.Empty{}
		}
	}
	for el := range cache {
		_, exists := elems[el]
		if !exists {
			rm[el] = types.Empty{}
		}
	}

	// Save
	d.catCache[strings.ToLower(name)] = elems
	if len(rm) > 0 {
		toRm := make([]int, len(rm))
		i := 0
		for k := range add {
			toRm[i] = k
			i++
		}
		entry := catCacheEntry{
			Op:   catCacheOpAdd,
			Data: toRm,
		}
		dat, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		_, err = file.WriteString(string(dat) + "\n")
		if err != nil {
			return err
		}
	}
	if len(add) > 0 {
		toAdd := make([]int, len(add))
		i := 0
		for k := range add {
			toAdd[i] = k
			i++
		}
		entry := catCacheEntry{
			Op:   catCacheOpAdd,
			Data: toAdd,
		}
		dat, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		_, err = file.WriteString(string(dat) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) DelCatCache(name string) error {
	d.Lock()
	defer d.Unlock()

	delete(d.catCache, strings.ToLower(name))
	delete(d.catCacheFiles, strings.ToLower(name))
	return os.Remove(filepath.Join(d.dbPath, "catcache", url.PathEscape(name)+".json"))
}

func (d *DB) SaveCat(elems *types.Category) error {
	// Empty?
	if len(elems.Elements) == 0 {
		d.Lock()
		delete(d.cats, strings.ToLower(elems.Name))
		delete(d.catFiles, strings.ToLower(elems.Name))
		d.Unlock()

		err := os.Remove(filepath.Join(d.dbPath, "categories", url.PathEscape(elems.Name)+".json"))
		if err != nil {
			return err
		}
		return d.DelCatCache(elems.Name)
	}

	elems.Lock.RLock()
	dat, err := json.Marshal(elems)
	elems.Lock.RUnlock()
	if err != nil {
		return err
	}

	d.Lock()

	file, exists := d.catFiles[strings.ToLower(elems.Name)]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "categories", url.PathEscape(elems.Name)+".json"))
		if err != nil {
			d.Unlock()
			return err
		}
		d.catFiles[strings.ToLower(elems.Name)] = file
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		d.Unlock()
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		d.Unlock()
		return err
	}
	_, err = file.Write(dat)
	if err != nil {
		d.Unlock()
		return err
	}
	d.Unlock()

	// Save cat cache
	elems.Lock.RLock()
	copy := make(map[int]types.Empty, len(elems.Elements))
	for k := range elems.Elements {
		copy[k] = types.Empty{}
	}
	elems.Lock.RUnlock()
	return d.SaveCatCache(elems.Name, copy)
}

func (d *DB) SaveInv(inv *types.Inventory, recalc ...bool) error {
	d.RLock()
	if len(recalc) == 1 {
		inv.MadeCnt = 0
		for elem := range inv.Elements {
			elem, res := d.GetElement(elem, true)
			if !res.Exists {
				continue
			}
			if elem.Creator == inv.User {
				inv.MadeCnt++
			}
		}
	} else if len(recalc) == 2 {
		inv.MadeCnt++
	}
	d.RUnlock()

	inv.Lock.RLock()
	dat, err := json.Marshal(inv)
	inv.Lock.RUnlock()
	if err != nil {
		return err
	}

	file, exists := d.invFiles[inv.User]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "inventories", inv.User+".json"))
		if err != nil {
			return err
		}
		d.invFiles[inv.User] = file
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
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
