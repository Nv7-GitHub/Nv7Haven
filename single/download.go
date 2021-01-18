package single

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Single) download(c *fiber.Ctx) error {
	id := c.Params("id")
	uid := c.Params("uid")
	return c.SendFile(fmt.Sprintf("/home/container/packs/%s_%s.pack", uid, id))
}
