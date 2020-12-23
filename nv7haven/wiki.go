package nv7haven

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) searchElems(c *fiber.Ctx) error {
	
	

	query, err := url.PathUnescape(c.Params("query"))
	if err != nil {
		return err
	}

	res, err := n.sql.Query("SELECT name FROM elements WHERE name LIKE ? LIMIT 100", query)
	if err != nil {
		return err
	}
	out := make([]string, 0)
	for res.Next() {
		var data string
		err = res.Scan(&data)
		if err != nil {
			return err
		}
		out = append(out, data)
	}

	return c.JSON(out)
}
