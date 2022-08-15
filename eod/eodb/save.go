package eodb

import (
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

	val := d.vcats[strings.ToLower(name)]
	if val.Rule == types.VirtualCategoryRuleRegex {
		d.Unlock()
		err := d.DelCatCache(val.Name)
		if err != nil {
			return err
		}
		d.Lock()
	}

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
	rm := make([]int, 0)
	add := make([]int, 0)
	for el := range elems {
		_, exists := cache[el]
		if !exists {
			add = append(add, el)
		}
	}
	for el := range cache {
		_, exists := elems[el]
		if !exists {
			rm = append(rm, el)
		}
	}

	// Copy cache
	v := make(map[int]types.Empty, len(elems))
	for k := range elems {
		v[k] = elems[k]
	}
	d.catCache[strings.ToLower(name)] = v

	// Save
	if len(rm) > 0 {
		entry := catCacheEntry{
			Op:   catCacheOpRemove,
			Data: rm,
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
		entry := catCacheEntry{
			Op:   catCacheOpAdd,
			Data: add,
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

type invOpKind int

const (
	invOpAdd invOpKind = 0
	invOpRem invOpKind = 1
)

type invOp struct {
	Kind invOpKind
	Data []int
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

		dataFile, err := os.Create(filepath.Join(d.dbPath, "invdata", inv.User+".json"))
		if err != nil {
			return err
		}
		d.invDataFiles[inv.User] = dataFile
		d.invData[inv.User] = make(map[int]types.Empty)
	}

	data := d.invData[inv.User]

	// Calc diff
	rm := make([]int, 0)
	add := make([]int, 0)
	inv.Lock.RLock()
	for el := range inv.Elements {
		_, exists := data[el]
		if !exists {
			add = append(add, el)
		}
	}
	for el := range data {
		_, exists := inv.Elements[el]
		if !exists {
			rm = append(rm, el)
		}
	}
	inv.Lock.RUnlock()

	// Write
	dataFile := d.invDataFiles[inv.User]
	if len(rm) > 0 {
		op := invOp{
			Kind: invOpRem,
			Data: rm,
		}
		dat, err := json.Marshal(op)
		if err != nil {
			return err
		}
		_, err = dataFile.WriteString(string(dat) + "\n")
		if err != nil {
			return err
		}
	}
	if len(add) > 0 {
		entry := invOp{
			Kind: invOpAdd,
			Data: add,
		}
		dat, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		_, err = dataFile.WriteString(string(dat) + "\n")
		if err != nil {
			return err
		}
	}

	// Save data
	inv.Lock.RLock()
	cop := make(map[int]types.Empty, len(inv.Elements))
	for k := range inv.Elements {
		cop[k] = types.Empty{}
	}
	inv.Lock.RUnlock()

	d.Lock()
	d.invData[inv.User] = cop
	d.Unlock()

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
