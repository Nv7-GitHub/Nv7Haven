package elemental

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) randomSuggestion(where string, uid string) ([]string, error) {
	isAnarchy := time.Now().Weekday() == anarchyDay
	params := []interface{}{maxVotes, "%\"" + uid + "\"%"}
	if isAnarchy {
		where = "1"
		params = []interface{}{}
	}

	res, err := e.db.Query("SELECT name FROM suggestions WHERE "+where+" LIMIT 100", params...)
	if err != nil {
		return []string{}, err
	}
	defer res.Close()
	comings := make([]string, 0)
	var name string
	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			return []string{}, err
		}
		comings = append(comings, name)
	}
	if len(comings) == 0 {
		return []string{}, nil
	}

	for tries := 0; tries < 10; tries++ {
		item := comings[rand.Intn(len(comings))]
		parents, err := e.getSuggParents(item)
		if err != nil {
			return []string{}, err
		}

		if len(parents) == 2 {
			if isAnarchy {
				return parents, nil
			}

			combos, err := e.getSuggNeighbors(parents[1])
			if err != nil {
				return []string{}, err
			}

			for _, sugg := range combos {
				sugg, err := e.getSugg(sugg)
				if err != nil {
					return []string{}, err
				}
				for _, val := range sugg.Voted {
					if val == uid {
						continue
					}
				}
			}

			return parents, nil
		}
	}
	return []string{}, nil
}

func (e *Elemental) upAndComingSuggestion(c *fiber.Ctx) error {
	ans, err := e.UpAndComingSuggestion(c.Params("uid"))
	if err != nil {
		return err
	}
	return c.JSON(ans)
}

// Pretty much the same, just different first line
func (e *Elemental) randomLonelySuggestion(c *fiber.Ctx) error {
	ans, err := e.RandomLonelySuggestion(c.Params("uid"))
	if err != nil {
		return err
	}
	return c.JSON(ans)
}

// RandomLonelySuggestion gets a random lonely suggestion
func (e *Elemental) RandomLonelySuggestion(uid string) ([]string, error) {
	return e.randomSuggestion("votes<? AND voted NOT LIKE ?", uid)
}

// UpAndComingSuggestion suggestion gets a suggestion that needs one vote
func (e *Elemental) UpAndComingSuggestion(uid string) ([]string, error) {
	return e.randomSuggestion("votes<? AND voted NOT LIKE ?", uid)
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
