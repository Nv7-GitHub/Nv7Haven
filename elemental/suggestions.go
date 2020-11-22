package elemental

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

const minVotes = -1

func getSugg(id string) (Suggestion, error) {
	data, err := db.Get("suggestions/" + id)
	if err != nil {
		return Suggestion{}, err
	}
	if string(data) == "null" {
		return Suggestion{}, errors.New("null")
	}
	var suggestion Suggestion
	err = json.Unmarshal(data, &suggestion)
	if err != nil {
		return Suggestion{}, err
	}
	return suggestion, nil
}

func getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	suggestion, err := getSugg(id)
	if err != nil {
		if err.Error() == "null" {
			return c.SendString("null")
		}
		return err
	}
	return c.JSON(suggestion)
}

func getSuggestionCombos(c *fiber.Ctx) error {
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
	comboData, err := db.Get("suggestionMap/" + elem1 + "/" + elem2)
	if err != nil {
		return err
	}
	var data []string
	err = json.Unmarshal(comboData, &data)
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func downVoteSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := getSugg(id)
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
		db.SetData("suggestions/"+id, nil)
	}
	existing.Voted = append(existing.Voted, uid)
	db.SetData("suggestion/"+id, existing)
	return nil
}
