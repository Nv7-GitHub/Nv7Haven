package nv7haven

import (
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type nameData struct {
	Name       string
	IsMale     bool
	Population int
}

func (n *Nv7Haven) searchNames(c *fiber.Ctx) error {
	query, err := url.PathUnescape(c.Params("query"))
	if err != nil {
		return err
	}

	res, err := n.sql.Query("SELECT name FROM names WHERE name LIKE ? ORDER BY count DESC LIMIT 100", query)
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

func (n *Nv7Haven) getName(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	row := n.sql.QueryRow("SELECT * FROM names WHERE name=? ORDER BY count DESC LIMIT 1", name)
	var nm nameData
	err = row.Scan(&nm.Name, &nm.IsMale, &nm.Population)
	if err != nil {
		return err
	}
	return c.JSON(nm)
}

func (n *Nv7Haven) nameCount(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	if name == "all" {
		row := n.sql.QueryRow("SELECT SUM(`count`) AS num FROM `names`")
		var num int
		err = row.Scan(&num)
		if err != nil {
			return err
		}
		c.WriteString(strconv.Itoa(num))
		return nil
	}

	row := n.sql.QueryRow("SELECT SUM(`count`) AS num FROM `names` WHERE name=?", name)
	var num int
	err = row.Scan(&num)
	if err != nil {
		return err
	}

	c.WriteString(strconv.Itoa(num))
	return nil
}
