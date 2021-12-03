package eodb

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	_, err := d.comboFile.WriteString(body + "=" + strconv.Itoa(result))
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

func (d *DB) SaveCat(elems *types.ElemContainer) error {
	elems.RLock()
	dat, err := json.Marshal(elems.Data)
	elems.RUnlock()
	if err != nil {
		return err
	}

	d.Lock()
	defer d.Unlock()

	file, exists := d.catFiles[strings.ToLower(elems.Id)]
	if !exists {
		file, err = os.Create(filepath.Join(d.dbPath, "categories", url.PathEscape(elems.Id)))
		if err != nil {
			return err
		}
		d.catFiles[strings.ToLower(elems.Id)] = file
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
		file, err = os.Create(filepath.Join(d.dbPath, "inventories", inv.Id))
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
