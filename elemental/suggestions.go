package elemental

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const minVotes = -1
const maxVotes = 3
const anarchyDay = 6

func (e *Elemental) getSugg(id string) (Suggestion, error) {
	res, err := e.db.Query("SELECT * FROM suggestions WHERE name=?", id)
	if err != nil {
		return Suggestion{}, err
	}
	defer res.Close()
	var suggestion Suggestion
	var color string
	var voted string
	res.Next()
	err = res.Scan(&suggestion.Name, &color, &suggestion.Creator, &voted, &suggestion.Votes)
	if err != nil {
		return Suggestion{}, err
	}

	colors := strings.Split(color, "_")
	sat, err := strconv.ParseFloat(colors[1], 32)
	if err != nil {
		return Suggestion{}, err
	}
	light, err := strconv.ParseFloat(colors[2], 32)
	if err != nil {
		return Suggestion{}, err
	}
	suggestion.Color = Color{
		Base:       colors[0],
		Saturation: float32(sat),
		Lightness:  float32(light),
	}

	var votedData []string
	err = json.Unmarshal([]byte(voted), &votedData)
	if err != nil {
		return Suggestion{}, err
	}
	suggestion.Voted = votedData

	return suggestion, nil
}

func (e *Elemental) getSuggestion(c *fiber.Ctx) error {
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	suggestion, err := e.getSugg(id)
	if err != nil {
		if err.Error() == "null" {
			return c.SendString("null")
		}
		return err
	}
	return c.JSON(suggestion)
}

func (e *Elemental) getSuggestionCombos(c *fiber.Ctx) error {
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}
	data, err := e.getSuggestions(elem1)
	if err != nil {
		return err
	}
	return c.JSON(data[elem2])
}

func (e *Elemental) downVoteSuggestion(c *fiber.Ctx) error {

	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := e.getSugg(id)
	if err != nil {
		return err
	}
	for _, voted := range existing.Voted {
		if voted == uid {
			return c.SendString("You already voted!")
		}
	}
	existing.Votes--
	if existing.Votes < minVotes {
		e.db.Exec("DELETE FROM suggestions WHERE name=?", id)
	}
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	if err != nil {
		return err
	}
	_, err = e.db.Exec("UPDATE suggestions SET voted=?, votes=? WHERE name=?", data, existing.Votes, existing.Name)
	if err != nil {
		return err
	}
	return nil
}

func (e *Elemental) upVoteSuggestion(c *fiber.Ctx) error {
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := e.getSugg(id)
	log.Println("err1", err)
	if err != nil {
		return err
	}

	isAnarchy := int(time.Now().Weekday()) == anarchyDay
	if !(isAnarchy) {
		for _, voted := range existing.Voted {
			if voted == uid {
				return c.SendString("You already voted!")
			}
		}
	}

	existing.Votes++
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	log.Println("err2", err)
	if err != nil {
		return err
	}
	_, err = e.db.Exec("UPDATE suggestions SET votes=?, voted=? WHERE name=?", existing.Votes, data, existing.Name)
	log.Println("err3", err)
	if err != nil {
		return err
	}
	if (existing.Votes >= maxVotes) || isAnarchy {
		return c.SendString("create")
	}
	return nil
}

func (e *Elemental) newSuggestion(c *fiber.Ctx) error {
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}
	newElem, err := url.PathUnescape(c.Params("data"))
	if err != nil {
		return err
	}

	var suggestion Suggestion
	err = json.Unmarshal([]byte(newElem), &suggestion)
	if err != nil {
		return err
	}

	voted, _ := json.Marshal(suggestion.Voted)
	color := fmt.Sprintf("%s_%f_%f", suggestion.Color.Base, suggestion.Color.Saturation, suggestion.Color.Lightness)
	_, err = e.db.Exec("INSERT INTO suggestions VALUES( ?, ?, ?, ?, ? )", suggestion.Name, color, suggestion.Creator, voted, suggestion.Votes)
	if err != nil {
		return err
	}

	combos, err := e.getSuggestions(elem1)
	if err != nil {
		return err
	}
	combos[elem2] = append(combos[elem2], suggestion.Name)
	data, err := json.Marshal(combos)
	if err != nil {
		return err
	}
	_, err = e.db.Exec("UPDATE suggestion_combos SET combos=? WHERE name=?", data, elem1)
	if err != nil {
		return err
	}
	if int(time.Now().Weekday()) == anarchyDay {
		return c.SendString("create")
	}
	return nil
}
