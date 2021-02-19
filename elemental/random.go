package elemental

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) upAndComingSuggestion(c *fiber.Ctx) error {
	isAnarchy := int(time.Now().Weekday()) == anarchyDay

	where := "votes=? AND voted NOT LIKE ?"
	params := []interface{}{maxVotes, "%\"" + c.Params("uid") + "\"%"}
	if isAnarchy {
		where = "1"
		params = []interface{}{}
	}

	res, err := e.db.Query("SELECT name FROM suggestions WHERE "+where+" LIMIT 100", params...)
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

	for tries := 0; tries < 10; tries++ {
		item := comings[rand.Intn(len(comings))]
		parents, err := e.getSuggParents(item)
		if err != nil {
			return err
		}

		if len(parents) == 2 {
			if isAnarchy {
				return c.JSON(parents)
			}

			data, err := e.getSuggestions(parents[0])
			if err != nil {
				return err
			}
			combos := data[parents[1]]

			for _, sugg := range combos {
				sugg, err := e.getSugg(sugg)
				if err != nil {
					return err
				}
				for _, val := range sugg.Voted {
					if val == c.Params("uid") {
						continue
					}
				}
			}

			return c.JSON(parents)
		}
	}
	return c.JSON([]string{})
}

// Pretty much the same, just different first line
func (e *Elemental) randomLonelySuggestion(c *fiber.Ctx) error {
	isAnarchy := int(time.Now().Weekday()) == anarchyDay

	where := "votes<? AND voted NOT LIKE ?"
	params := []interface{}{maxVotes - 1, "%\"" + c.Params("uid") + "\"%"}
	if isAnarchy {
		where = "1"
		params = []interface{}{}
	}

	res, err := e.db.Query("SELECT name FROM suggestions WHERE "+where+" LIMIT 100", params...)
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

	for tries := 0; tries < 10; tries++ {
		item := comings[rand.Intn(len(comings))]
		parents, err := e.getSuggParents(item)
		if err != nil {
			return err
		}

		if len(parents) == 2 {
			if isAnarchy {
				return c.JSON(parents)
			}

			data, err := e.getSuggestions(parents[0])
			if err != nil {
				return err
			}
			combos := data[parents[1]]

			for _, sugg := range combos {
				sugg, err := e.getSugg(sugg)
				if err != nil {
					return err
				}
				for _, val := range sugg.Voted {
					if val == c.Params("uid") {
						continue
					}
				}
			}

			return c.JSON(parents)
		}
	}
	return c.JSON([]string{})
}

func (e *Elemental) getSuggParents(item string) ([]string, error) {
	res, err := e.db.Query("SELECT * FROM suggestion_combos WHERE combos LIKE ?", "%\""+item+"\"%")
	if err != nil {
		return nil, err
	}
	var name string
	var combos string
	var comboDat map[string][]string
	for res.Next() {
		err = res.Scan(&name, &combos)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(combos), &comboDat)
		if err != nil {
			return nil, err
		}
		for k, v := range comboDat {
			for _, val := range v {
				if val == item {
					return []string{name, k}, nil
				}
			}
		}
	}
	return []string{}, nil
}
