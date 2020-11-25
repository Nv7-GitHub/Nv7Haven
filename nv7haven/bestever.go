package nv7haven

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

func getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	min := len(data) - 50
	if min < 0 {
		min = 0
	}
	randNum1 := rand.Intn(len(data)-min) + min
	randNum2 := rand.Intn(len(data)-min) + min
	for randNum2 == randNum1 {
		randNum2 = rand.Intn(len(data)-min) + min
	}
	output := map[int]string{
		randNum1: data[randNum1].Name,
		randNum2: data[randNum2].Name,
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
	end := len(data) - 1
	if end > 9 {
		end = 9
	}
	dat := make([]string, end+1)
	i := 0
	for _, val := range data[:end+1] {
		dat[i] = val.Name
		i++
	}
	return c.JSON(dat)
}
