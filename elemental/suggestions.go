package elemental

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const minVotes = -1
const maxVotes = 3
const anarchyDay = time.Saturday

func (e *Elemental) getSugg(id string) (Suggestion, error) {
	row := e.db.QueryRow("SELECT * FROM suggestions WHERE name=?", id)
	var suggestion Suggestion
	var color string
	var voted string
	err := row.Scan(&suggestion.Name, &color, &suggestion.Creator, &voted, &suggestion.Votes)
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
	data, err := e.GetSuggestions(elem1, elem2)
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func (e *Elemental) downVoteSuggestion(c *fiber.Ctx) error {
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	suc, msg := e.DownvoteSuggestion(id, uid)
	if !suc {
		return errors.New(msg)
	}

	return nil
}

// DownvoteSuggestion downvotes a suggestion
func (e *Elemental) DownvoteSuggestion(id, uid string) (bool, string) {
	existing, err := e.getSugg(id)
	if err != nil {
		return false, err.Error()
	}
	for _, voted := range existing.Voted {
		if voted == uid {
			return false, "You already voted!"
		}
	}
	existing.Votes--
	if existing.Votes < minVotes {
		e.db.Exec("DELETE FROM suggestions WHERE name=?", id)
		return true, ""
	}
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	if err != nil {
		return false, err.Error()
	}
	_, err = e.db.Exec("UPDATE suggestions SET voted=?, votes=? WHERE name=?", data, existing.Votes, existing.Name)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func (e *Elemental) upVoteSuggestion(c *fiber.Ctx) error {
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	create, suc, msg := e.UpvoteSuggestion(id, uid)
	if !suc {
		return errors.New(msg)
	}

	if create {
		return c.SendString("create")
	}

	return nil
}

// UpvoteSuggestion upvotes a suggestion
func (e *Elemental) UpvoteSuggestion(id, uid string) (bool, bool, string) {
	existing, err := e.getSugg(id)
	if err != nil {
		return false, false, err.Error()
	}

	isAnarchy := time.Now().Weekday() == anarchyDay
	if !(isAnarchy) {
		for _, voted := range existing.Voted {
			if voted == uid {
				return false, false, "You already voted!"
			}
		}
	}

	existing.Votes++
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	if err != nil {
		return false, false, err.Error()
	}
	_, err = e.db.Exec("UPDATE suggestions SET votes=?, voted=? WHERE name=?", existing.Votes, data, existing.Name)
	if err != nil {
		return false, false, err.Error()
	}
	if (existing.Votes >= maxVotes) || isAnarchy {
		return true, true, ""
	}
	return false, true, ""
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

	create, err := e.NewSuggestion(elem1, elem2, suggestion)
	if err != nil {
		return err
	}

	if create {
		return c.SendString("create")
	}

	return nil
}

// NewSuggestion makes a new suggestion
func (e *Elemental) NewSuggestion(elem1, elem2 string, suggestion Suggestion) (bool, error) {
	voted, _ := json.Marshal(suggestion.Voted)
	color := fmt.Sprintf("%s_%f_%f", suggestion.Color.Base, suggestion.Color.Saturation, suggestion.Color.Lightness)
	_, err := e.db.Exec("INSERT INTO suggestions VALUES( ?, ?, ?, ?, ? )", suggestion.Name, color, suggestion.Creator, voted, suggestion.Votes)
	if err != nil {
		return false, err
	}

	_, err = e.db.Exec("INSERT INTO sugg_combos VALUES ( ?, ?, ? )", elem1, elem2, suggestion.Name)
	if err != nil {
		return false, err
	}
	if time.Now().Weekday() == anarchyDay {
		return true, nil
	}
	return false, nil
}
