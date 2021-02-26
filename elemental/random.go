package elemental

import (
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

			combos, err := e.getSuggNeighbors(parents[1])
			if err != nil {
				return err
			}

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

			combos, err := e.getSuggNeighbors(parents[1])
			if err != nil {
				return err
			}

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
	row := e.db.QueryRow("SELECT elem1, elem2 FROM sugg_combos WHERE elem3=?", item)
	var elem1 string
	var elem2 string
	err := row.Scan(&elem1, &elem2)
	if err != nil {
		return nil, err
	}
	return []string{elem1, elem2}, nil
}

func (e *Elemental) getSuggNeighbors(elem1 string) ([]string, error) {
	res, err := e.db.Query("SELECT elem3 FROM sugg_combos WHERE elem1=? OR elem2=?", elem1, elem1)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var dat string
	var out []string
	for res.Next() {
		err = res.Scan(&dat)
		if err != nil {
			return nil, err
		}
		out = append(out, dat)
	}
	return out, nil
}
