package elemental

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Element has the data for a created element
type Element struct {
	Color     string   `json:"color"`
	Comment   string   `json:"comment"`
	CreatedOn int      `json:"createdOn"`
	Creator   string   `json:"creator"`
	Name      string   `json:"name"`
	Parents   []string `json:"parents"`
	Pioneer   string   `json:"pioneer"`
	Uses      int      `json:"uses"`
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

func (e *Elemental) getElement(elemName string) (Element, error) {
	val, exists := e.cache[elemName]
	if !exists {
		var elem Element
		res, err := e.db.Query("SELECT * FROM elements WHERE name=?", elemName)
		if err != nil {
			return Element{}, err
		}
		defer res.Close()
		elem.Parents = make([]string, 2)
		res.Next()
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn)
		if err != nil {
			return Element{}, err
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}

		fmt.Println(elem.Name)
		uses, err := e.db.Query("SELECT COUNT(1) FROM elem_combos WHERE elem1=? OR elem2=?", elem.Name, elem.Name)
		if err != nil {
			return Element{}, err
		}
		uses.Next()
		err = uses.Scan(&elem.Uses)
		if err != nil {
			return Element{}, err
		}

		var mutex = &sync.RWMutex{}
		mutex.Lock()
		e.cache[elemName] = elem
		mutex.Unlock()
		return elem, nil
	}
	return val, nil
}

func (e *Elemental) getElem(c *fiber.Ctx) error {
	elemName, err := url.PathUnescape(c.Params("elem"))
	if err != nil {
		return err
	}
	elem, err := e.getElement(elemName)
	if err != nil {
		return err
	}
	return c.JSON(elem)
}

func (e *Elemental) getCombo(c *fiber.Ctx) error {
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}

	res, err := e.db.Query("SELECT COUNT(1) FROM elem_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?) LIMIT 1", elem1, elem2, elem2, elem1)
	if err != nil {
		return err
	}
	defer res.Close()
	var count int
	res.Next()
	res.Scan(&count)
	if count == 0 {
		return c.JSON(map[string]bool{
			"exists": false,
		})
	}

	res, err = e.db.Query("SELECT elem3 FROM elem_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?) LIMIT 1", elem1, elem2, elem2, elem1)
	if err != nil {
		return err
	}
	defer res.Close()
	var elem3 string
	res.Next()
	err = res.Scan(&elem3)
	if err != nil {
		return err
	}

	return c.JSON(map[string]interface{}{
		"exists": true,
		"combo":  elem3,
	})
}

func (e *Elemental) getAll(c *fiber.Ctx) error {
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", c.Params("uid"))
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	var found []string
	err = json.Unmarshal([]byte(data), &found)
	if err != nil {
		return err
	}

	var recents []RecentCombination
	dat, err := e.fdb.Get("recent")
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &recents)
	if err != nil {
		return err
	}

	req := make(map[string]bool)
	for _, val := range found {
		req[val] = true
	}
	for _, rec := range recents {
		req[rec.Result] = true
		req[rec.Recipe[0]] = true
		req[rec.Recipe[1]] = true
	}

	final := make([]Element, len(req))
	i := 0
	for k := range req {
		final[i], err = e.getElement(k)
		if err != nil {
			return err
		}
		i++
	}
	return c.JSON(final)
}
