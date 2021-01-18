package single

import (
	"fmt"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
)

func (s *Single) download(c *fiber.Ctx) error {
	id := c.Params("id")
	uid := c.Params("uid")
	dat, err := ioutil.ReadFile(fmt.Sprintf("/home/container/packs/%s_%s.pack", uid, id))
	if err != nil {
		return err
	}
	return c.SendString(string(dat))
}
