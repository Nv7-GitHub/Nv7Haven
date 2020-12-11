package elemental

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

const minVotes = -1
const maxVotes = 3 // ANARCHY: 0, ORIGINAL: 3

func (e *Elemental) getSugg(id string) (Suggestion, error) {
	data, err := e.db.Get("suggestions/" + url.PathEscape(id))
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

func (e *Elemental) getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
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
	comboData, err := e.db.Get("suggestionMap/" + url.PathEscape(elem1) + "/" + url.PathEscape(elem2))
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

func (e *Elemental) downVoteSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
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
		e.db.SetData("suggestions/"+url.PathEscape(id), nil)
	}
	existing.Voted = append(existing.Voted, uid)
	err = e.db.SetData("suggestions/"+url.PathEscape(id), existing)
	if err != nil {
		return err
	}
	return nil
}

func (e *Elemental) upVoteSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := e.getSugg(id)
	if err != nil {
		return err
	}
	// ANARCHY: Comment out this section
	for _, voted := range existing.Voted {
		if voted == uid {
			return c.SendString("You already voted!")
		}
	}
	// ANARCHY
	existing.Votes++
	existing.Voted = append(existing.Voted, uid)
	err = e.db.SetData("suggestions/"+url.PathEscape(id), existing)
	if err != nil {
		return err
	}
	if existing.Votes >= maxVotes {
		return c.SendString("create")
	}
	return nil
}

func (e *Elemental) newSuggestion(c *fiber.Ctx) error {
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
	newElem, err := url.PathUnescape(c.Params("data"))
	if err != nil {
		return err
	}

	var suggestion Suggestion
	err = json.Unmarshal([]byte(newElem), &suggestion)
	if err != nil {
		return err
	}

	err = e.db.SetData("suggestions/"+url.PathEscape(suggestion.Name), suggestion)
	if err != nil {
		return err
	}

	comboData, err := e.db.Get("suggestionMap/" + url.PathEscape(elem1) + "/" + url.PathEscape(elem2))
	if err != nil {
		return err
	}
	var data []string
	err = json.Unmarshal(comboData, &data)
	if err != nil {
		return err
	}
	data = append(data, suggestion.Name)
	e.db.SetData("suggestionMap/"+url.PathEscape(elem1)+"/"+url.PathEscape(elem2), data)
	return nil
}
