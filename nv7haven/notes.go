package nv7haven

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) newNote(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return err
	}
	ip := c.IP()

	// Does it exist?
	res, err := n.sql.Query("SELECT COUNT(*) FROM notes WHERE ip=? AND name=?", ip, name)
	if err != nil {
		return err
	}
	defer res.Close()
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

func (n *Nv7Haven) changeNote(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return err
	}
	ip := c.IP()

	body := string(c.Body())

	// Create note
	_, err = n.sql.Exec("UPDATE notes SET note=? WHERE name=? AND password=? AND ip=?", body, name, password, ip)
	if err != nil {
		return err
	}

	return nil
}

func (n *Nv7Haven) getNote(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return err
	}
	ip := c.IP()

	res, err := n.sql.Query("SELECT note FROM notes WHERE ip=? AND name=? AND password=?", ip, name, password)
	if err != nil {
		return err
	}
	defer res.Close()
	res.Next()
	var data string
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	return c.SendString(data)
}

func (n *Nv7Haven) hasPassword(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	ip := c.IP()

	res, err := n.sql.Query("SELECT password FROM notes WHERE ip=? AND name=?", ip, name)
	if err != nil {
		return err
	}
	defer res.Close()
	res.Next()
	var data string
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	if data != "" {
		return c.SendString("1")
	}
	return c.SendString("0")
}

func (n *Nv7Haven) searchNotes(c *fiber.Ctx) error {
	query, err := url.PathUnescape(c.Params("query"))
	if err != nil {
		return err
	}
	ip := c.IP()

	res, err := n.sql.Query("SELECT name FROM notes WHERE ip=? AND name LIKE ?", ip, query)
	if err != nil {
		return err
	}
	defer res.Close()
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

func (n *Nv7Haven) deleteNote(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return err
	}
	ip := c.IP()

	_, err = n.sql.Exec("DELETE FROM notes WHERE ip=? AND name=? AND password=?", ip, name, password)
	if err != nil {
		return err
	}
	return nil
}
