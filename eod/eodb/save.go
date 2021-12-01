package eodb

import (
	"encoding/json"
	"strconv"

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
