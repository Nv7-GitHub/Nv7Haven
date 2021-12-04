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
		el.ID = len(d.elements) + 1
		d.elements = append(d.elements, el)
		d.elemNames[strings.ToLower(el.Name)] = el.ID
	} else {
		d.elements[el.ID-1] = el
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

func (d *DB) SaveServerConfig() error {
	d.Lock()
	defer d.Unlock()

	dat, err := json.Marshal(d.config)
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
	elems.Lock.RLock()
	dat, err := json.Marshal(elems.Elements)
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
	_, err = file.Write(dat)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SaveInv(inv *types.ElemContainer) error {
	inv.RLock()
	dat, err := json.Marshal(inv)
	inv.RUnlock()
	if err != nil {
		return err
	}

	file, exists := d.invFiles[inv.Id]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "inventories", inv.Id+".json"))
		if err != nil {
			return err
		}
		d.invFiles[inv.Id] = file
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
