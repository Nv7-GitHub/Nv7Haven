package elemental

import (
	"encoding/json"
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
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

func (e *Elemental) getElem(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	elemName, err := url.PathUnescape(c.Params("elem"))
	if err != nil {
		return err
	}
	val, exists := e.cache[elemName]
	if !exists {
		var elem Element
		res, err := e.db.Query("SELECT * FROM elements WHERE name=?", elemName)
		if err != nil {
			return err
		}
		defer res.Close()
		elem.Parents = make([]string, 2)
		res.Next()
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn)
		if err != nil {
			return err
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		var mutex = &sync.RWMutex{}
		mutex.Lock()
		e.cache[elemName] = elem
		mutex.Unlock()
		return c.JSON(elem)
	}
	return c.JSON(val)
}

func (e *Elemental) getCombo(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}

	res, err := e.db.Query("SELECT COUNT(1) FROM element_combos WHERE name=? LIMIT 1", elem1)
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

	var data map[string]string
	res, err = e.db.Query("SELECT combos FROM element_combos WHERE name=? LIMIT 1", elem1)
	if err != nil {
		return err
	}
	defer res.Close()
	var comboData string
	res.Next()
	err = res.Scan(&comboData)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(comboData), &data)
	if err != nil {
		return err
	}

	output, exists := data[elem2]
	if !exists {
		return c.JSON(map[string]bool{
			"exists": false,
		})
	}

	return c.JSON(map[string]interface{}{
		"exists": true,
		"combo":  output,
	})
}
