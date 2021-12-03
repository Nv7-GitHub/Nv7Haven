package eodb

import (
	"bufio"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (d *DB) loadElements() error {
	f, err := os.OpenFile(filepath.Join(d.dbPath, "elements.json"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)

	dat := types.Element{}
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// Parse
		err = json.Unmarshal(line, &dat)
		if err != nil {
			return err
		}

		// Add to elements
		if len(d.elements) < dat.ID {
			d.elements = append(d.elements, make([]types.Element, dat.ID-len(d.elements))...) // Grow
			d.elements[dat.ID-1] = dat
			d.elemNames[strings.ToLower(dat.Name)] = dat.ID
			dat = types.Element{}
		}
	}

	// Save
	d.elemFile = f
	return nil
}
func (d *DB) loadCombos() error {
	f, err := os.OpenFile(filepath.Join(d.dbPath, "combos.txt"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// Parse
		parts := strings.Split(string(line), "=") // 1+1=5
		result, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		d.combos[parts[0]] = result
	}

	// Save
	d.comboFile = f
	return nil
}

func (d *DB) loadConfig() error {
	f, err := os.OpenFile(filepath.Join(d.dbPath, "config.json"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	dat, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &d.config)
	if err != nil {
		return err
	}
	d.configFile = f

	return nil
}

func (d *DB) loadInvs() error {
	err := os.MkdirAll(filepath.Join(d.dbPath, "inventories"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.dbPath, "inventories"))
	if err != nil {
		return err
	}

	inv := make(map[int]types.Empty)
	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".json")
		f, err := os.Open(filepath.Join(d.dbPath, "inventories", file.Name()))
		if err != nil {
			return err
		}

		// Read inv
		dat, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		err = json.Unmarshal(dat, &inv)
		if err != nil {
			return err
		}

		// Save inv
		d.invs[name] = types.NewElemContainer(inv, name)
		d.invFiles[name] = f
		inv = make(map[int]types.Empty)
	}
	return nil
}

func (d *DB) loadCats() error {
	err := os.MkdirAll(filepath.Join(d.dbPath, "categories"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.dbPath, "categories"))
	if err != nil {
		return err
	}

	cat := make(map[int]types.Empty)
	for _, file := range files {
		name, err := url.PathUnescape(strings.TrimSuffix(file.Name(), ".json"))
		if err != nil {
			return err
		}
		f, err := os.Open(filepath.Join(d.dbPath, "categories", file.Name()))
		if err != nil {
			return err
		}

		// Read cat
		dat, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		err = json.Unmarshal(dat, &cat)
		if err != nil {
			return err
		}

		// Save cat
		d.cats[strings.ToLower(name)] = types.NewElemContainer(cat, name)
		d.catFiles[strings.ToLower(name)] = f
		cat = make(map[int]types.Empty)
	}
	return nil
}

func (d *DB) loadPolls() error {
	err := os.MkdirAll(filepath.Join(d.dbPath, "polls"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.dbPath, "polls"))
	if err != nil {
		return err
	}

	poll := types.Poll{}
	for _, file := range files {
		f, err := os.Open(filepath.Join(d.dbPath, "polls", file.Name()))
		if err != nil {
			return err
		}

		// Read poll
		dat, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		err = json.Unmarshal(dat, &poll)
		if err != nil {
			return err
		}

		// Save poll
		d.polls[poll.Message] = poll
		f.Close()
	}
	return nil
}
