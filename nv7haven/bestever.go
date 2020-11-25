package nv7haven

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmcvetta/randutil"
)

// Suggestion represents the datatype for a suggestion
type Suggestion struct {
	Votes int
	Name  string
}

var data []Suggestion
var changes int

const required = 3

func changed() error {
	changes++
	if changes > required {
		err := db.SetData("", map[string][]Suggestion{
			"data": data,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func initBestEver() error {
	rand.Seed(time.Now().UnixNano())
	rawData, err := db.Get("")
	if err != nil {
		return err
	}
	var rawMarshaled map[string][]Suggestion
	err = json.Unmarshal(rawData, &rawMarshaled)
	if err != nil {
		return err
	}
	data = rawMarshaled["data"]
	return nil
}

func newSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	suggest, err := url.PathUnescape(c.Params("suggestion"))
	if err != nil {
		return err
	}
	suggestion := Suggestion{
		Votes: 0,
		Name:  suggest,
	}
	for _, val := range data {
		if val.Name == suggest {
			return nil
		}
	}
	data = append(data, suggestion)
	changes = required
	err = changed()
	if err != nil {
		return err
	}
	return nil
}

type itemData struct {
	Name  string
	Index int
}

func getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	dat := make([]randutil.Choice, len(data))
	for i, val := range data {
		dat[i] = randutil.Choice{
			Weight: i + 1,
			Item: itemData{
				Name:  val.Name,
				Index: i,
			},
		}
	}
	choice1, err := randutil.WeightedChoice(dat)
	if err != nil {
		return err
	}
	choice2, err := randutil.WeightedChoice(dat)
	if err != nil {
		return err
	}
	for choice2.Item.(itemData).Index == choice1.Item.(itemData).Index {
		choice2, err = randutil.WeightedChoice(dat)
		if err != nil {
			return err
		}
	}
	output := map[int]string{
		choice1.Item.(itemData).Index: choice1.Item.(itemData).Name,
		choice2.Item.(itemData).Index: choice2.Item.(itemData).Name,
	}
	return c.JSON(output)
}

func vote(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	item, err := strconv.Atoi(c.Params("item"))
	if err != nil {
		return err
	}
	data[item].Votes++
	for !(item <= 0) && (data[item].Votes > data[item-1].Votes) {
		buffer := data[item-1]
		data[item-1] = data[item]
		data[item] = buffer
		item--
		changes = required
	}
	return changed()
}

func getLdb(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	end, err := strconv.Atoi(c.Params("len"))
	if err != nil {
		return c.JSON([]string{"Invalid input", "error: " + err.Error()})
	}
	end--
	if end > len(data)-1 {
		end = len(data) - 1
	}
	dat := make([]string, end+1)
	i := 0
	for _, val := range data[:end+1] {
		dat[i] = val.Name
		i++
	}
	return c.JSON(dat)
}
func refresh(c *fiber.Ctx) error {
	rawData, err := db.Get("")
	if err != nil {
		return err
	}
	var rawMarshaled map[string][]Suggestion
	err = json.Unmarshal(rawData, &rawMarshaled)
	if err != nil {
		return err
	}
	data = rawMarshaled["data"]
	return c.SendString("Success")
}
