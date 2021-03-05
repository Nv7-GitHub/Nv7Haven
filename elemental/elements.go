package elemental

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Element has the data for a created element
type Element struct {
	Color      string   `json:"color"`
	Comment    string   `json:"comment"`
	CreatedOn  int      `json:"createdOn"`
	Creator    string   `json:"creator"`
	Name       string   `json:"name"`
	Parents    []string `json:"parents"`
	Pioneer    string   `json:"pioneer"`
	Uses       int      `json:"uses"`
	FoundBy    int      `json:"foundby"`
	Complexity int      `json:"complexity"`
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

var lock = &sync.RWMutex{}

func (e *Elemental) calcComplexity(elem Element) (int, error) {
	if len(elem.Parents) == 0 {
		return 0, nil
	}
	parent1, err := e.GetElement(elem.Parents[0])
	if err != nil {
		return 0, err
	}
	parent2, err := e.GetElement(elem.Parents[1])
	if err != nil {
		return 0, err
	}
	comp1, err := e.calcComplexity(parent1)
	if err != nil {
		return 0, err
	}
	comp2, err := e.calcComplexity(parent2)
	if err != nil {
		return 0, err
	}
	return max(comp1, comp2) + 1, nil
}

// GetElement gets an element from the database
func (e *Elemental) GetElement(elemName string) (Element, error) {
	lock.RLock()
	val, exists := e.cache[elemName]
	lock.RUnlock()
	if !exists {
		var elem Element
		res, err := e.db.Query("SELECT * FROM elements WHERE name=?", elemName)
		if err != nil {
			return Element{}, err
		}
		defer res.Close()
		elem.Parents = make([]string, 2)
		res.Next()
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn, &elem.Complexity, &elem.Uses, &elem.FoundBy)
		if err != nil {
			return Element{}, err
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}

		lock.Lock()
		e.cache[elemName] = elem
		lock.Unlock()
		return elem, nil
	}
	return val, nil
}

func (e *Elemental) getElem(c *fiber.Ctx) error {
	elemName, err := url.PathUnescape(c.Params("elem"))
	if err != nil {
		return err
	}
	elem, err := e.GetElement(elemName)
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

	comb, suc, err := e.GetCombo(elem1, elem2)
	if err != nil {
		return err
	}
	if !suc {
		return c.JSON(map[string]interface{}{
			"exists": false,
		})
	}

	return c.JSON(map[string]interface{}{
		"exists": true,
		"combo":  comb,
	})
}

// GetCombo gets a combination
func (e *Elemental) GetCombo(elem1, elem2 string) (string, bool, error) {
	el1 := strings.ToUpper(elem1)
	el2 := strings.ToUpper(elem2)
	res, err := e.db.Query("SELECT COUNT(1) FROM elem_combos WHERE (UPPER(elem1)=? AND UPPER(elem2)=?) OR (UPPER(elem1)=? AND UPPER(elem2)=?) LIMIT 1", el1, el2, el2, el1)
	if err != nil {
		return "", false, err
	}
	defer res.Close()
	var count int
	res.Next()
	res.Scan(&count)
	if count == 0 {
		return "", false, nil
	}

	res, err = e.db.Query("SELECT elem3 FROM elem_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?) LIMIT 1", elem1, elem2, elem2, elem1)
	if err != nil {
		return "", false, err
	}
	defer res.Close()
	var elem3 string
	res.Next()
	err = res.Scan(&elem3)
	if err != nil {
		return "", false, err
	}

	return elem3, true, nil
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
		final[i], err = e.GetElement(k)
		if err != nil {
			return err
		}
		i++
	}
	return c.JSON(final)
}
