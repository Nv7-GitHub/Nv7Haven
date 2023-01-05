package nv7haven

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const pageLength = 50

var ldbQueryMap = map[string]string{
	"player": "SELECT name, JSON_LENGTH(found) AS found FROM `users` ORDER BY JSON_LENGTH(found) %s LIMIT ? OFFSET ?",
	"color":  `SELECT a.col AS col, (SELECT COUNT(*) AS cnt FROM elements WHERE SUBSTRING_INDEX(color, "_", 1)=a.col) FROM (SELECT DISTINCT SUBSTRING_INDEX(elements.color, "_", 1) AS col FROM elements) a ORDER BY (SELECT COUNT(*) FROM elements WHERE SUBSTRING_INDEX(color, "_", 1)=a.col) %s LIMIT ? OFFSET ? `,
}

type ldbReturn struct {
	PageLength int
	Items      []ldbItem
}

type ldbItem struct {
	Title string
	Value int
}

func (n *Nv7Haven) ldbQuery(c *fiber.Ctx) error {
	kind := c.Params("kind")
	order := "DESC"
	ord := c.Params("order")
	if ord == "1" {
		order = "ASC"
	}
	query, exists := ldbQueryMap[kind]
	if !exists {
		c.SendString("Invalid query type!")
		return fiber.ErrNotFound
	}
	page, err := strconv.Atoi(c.Params("page"))
	if err != nil {
		return err
	}
	res, err := n.sql.Query(fmt.Sprintf(query, order), pageLength, page*pageLength)
	if err != nil {
		return err
	}
	defer res.Close()
	out := ldbReturn{
		PageLength: pageLength,
		Items:      make([]ldbItem, 0),
	}
	for res.Next() {
		item := ldbItem{}
		err = res.Scan(&item.Title, &item.Value)
		if err != nil {
			return err
		}
		out.Items = append(out.Items, item)
	}
	return c.JSON(out)
}
