package elemental

import (
	"context"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
)

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
		var mutex = &sync.Mutex{}
		mutex.Lock()
		cache[elemName] = elem
		mutex.Unlock()
		return c.JSON(elem)
	}
	return c.JSON(val)
}
