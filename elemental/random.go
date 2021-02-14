package elemental

import (
	"encoding/json"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) upAndComingSuggestion(c *fiber.Ctx) error {
	res, err := e.db.Query("SELECT name FROM suggestions WHERE votes=? LIMIT 100 AND voted NOT LIKE ?", maxVotes, "%\""+c.Params("uid")+"\"%")
	if err != nil {
		return err
	}
	comings := make([]string, 0)
	var name string
	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			return err
		}
		comings = append(comings, name)
	}
	if len(comings) == 0 {
		return c.JSON([]string{})
	}
	item := comings[rand.Intn(len(comings))]

	res, err = e.db.Query("SELECT * FROM suggestion_combos WHERE combos LIKE ?", "%\""+item+"\"%")
	if err != nil {
		return err
	}
	var combos string
	var comboDat map[string][]string
	for res.Next() {
		err = res.Scan(&name, &combos)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(combos), &comboDat)
		if err != nil {
			return err
		}
		for k, v := range comboDat {
			for _, val := range v {
				if val == item {
					return c.JSON([]string{name, k})
				}
			}
		}
	}
	return c.JSON([]string{})
}

// Pretty much the same, just different first line
func (e *Elemental) randomLonelySuggestion(c *fiber.Ctx) error {
	res, err := e.db.Query("SELECT name FROM suggestions WHERE votes<? LIMIT 100 AND voted NOT LIKE ?", maxVotes-1, "%\""+c.Params("uid")+"\"%")
	if err != nil {
		return err
	}
	comings := make([]string, 0)
	var name string
	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			return err
		}
		comings = append(comings, name)
	}
	if len(comings) == 0 {
		return c.JSON([]string{})
	}
	item := comings[rand.Intn(len(comings))]

	res, err = e.db.Query("SELECT * FROM suggestion_combos WHERE combos LIKE ?", "%\""+item+"\"%")
	if err != nil {
		return err
	}
	var combos string
	var comboDat map[string][]string
	for res.Next() {
		err = res.Scan(&name, &combos)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(combos), &comboDat)
		if err != nil {
			return err
		}
		for k, v := range comboDat {
			for _, val := range v {
				if val == item {
					return c.JSON([]string{name, k})
				}
			}
		}
	}
	return c.JSON([]string{})
}
