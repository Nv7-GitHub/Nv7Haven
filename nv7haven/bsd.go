package nv7haven

import (
	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) bsdSearch(c *fiber.Ctx) error {
	query := "%" + c.Query("q") + "%"
	if query == "%%" {
		return c.SendStatus(400)
	}

	data := make([]map[string]any, 0)
	rows, err := n.bsdsql.Queryx("SELECT * FROM grading_data WHERE class_name LIKE ? OR teacher LIKE ? LIMIT 30", query, query)
	if err != nil {
		return err
	}
	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			return err
		}
		data = append(data, result)
	}

	return c.JSON(data)
}

func (n *Nv7Haven) bsdData(c *fiber.Ctx) error {
	id := c.Params("id")
	data := make(map[string]any)
	row := n.bsdsql.QueryRowx("SELECT * FROM grading_data WHERE id=?", id)
	err := row.MapScan(data)
	if err != nil {
		return err
	}

	return c.JSON(data)
}
