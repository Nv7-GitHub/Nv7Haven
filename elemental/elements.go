package elemental

import (
	"context"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Element has the data for a created element
type Element struct {
	Color     string   `json:"color"`
	Comment   string   `json:"comment"`
	CreatedOn int      `json:"createdOn"`
	Creator   string   `json:"creator"`
	Name      string   `json:"name"`
	Parents   []string `json:"parents"`
	Pioneer   string   `json:"pioneer"`
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

var cache map[string]Element = make(map[string]Element, 0)

func getElem(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	ctx := context.Background()
	elemName, err := url.PathUnescape(c.Params("elem"))
	if err != nil {
		return err
	}
	val, exists := cache[elemName]
	if !exists {
		var elem Element
		data, err := store.Collection("elements").Doc(elemName).Get(ctx)
		if err != nil {
			return err
		}
		err = data.DataTo(&elem)
		if err != nil {
			return err
		}
		var mutex = &sync.RWMutex{}
		mutex.Lock()
		cache[elemName] = elem
		mutex.Unlock()
		return c.JSON(elem)
	}
	return c.JSON(val)
}

func getCombo(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	ctx := context.Background()
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}
	var data map[string]string
	snapshot, err := store.Collection("combos").Doc(elem1).Get(ctx)
	if snapshot == nil || (snapshot.Exists() && err != nil) {
		return err
	} else if !snapshot.Exists() {
		return c.JSON(map[string]bool{
			"exists": false,
		})
	}
	err = snapshot.DataTo(&data)
	if err != nil {
		return err
	}
	output, exists := data[elem2]
	if !exists {
		return c.JSON(map[string]bool{
			"exists": false,
		})
	}
	return c.JSON(map[string]interface{}{
		"exists": true,
		"combo":  output,
	})
}
