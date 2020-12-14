package nv7haven

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) newNote(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return err
	}
	ip := c.IPs()[0]

	// Does it exist?
	res, err := n.sql.Query("SELECT COUNT(1) FROM notes WHERE ip=? AND name=?", ip, name)
	defer res.Close()
	if err != nil {
		return err
	}
	var count int
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return err
	}
	if count != 0 {
		return c.SendString("Note already exists. Try another name?")
	}

	// Create note
	_, err = n.sql.Exec("INSERT INTO notes VALUES ( ?, ?, ?, ? )", ip, name, password, "")
	if err != nil {
		return err
	}

	return nil
}
