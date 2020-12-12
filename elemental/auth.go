package elemental

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) createUser(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}

	uid, err := GenerateRandomStringURLSafe(16)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}

	// Check if name taken
	res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=\"?\" LIMIT 1", name, password)
	if err != nil {
		return err
	}
	defer res.Close()
	var count int
	err = res.Scan(&count)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	if count == 1 {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    "Account already exists!",
		})
	}

	_, err = e.db.Exec("INSERT INTO users VALUES( ?, ?, ?, ? )", name, uid, password, `["Air", "Earth", "Fire", "Water"]`)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}

	return c.JSON(map[string]interface{}{
		"success": true,
		"data":    uid,
	})
}

func (e *Elemental) loginUser(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	password, err := url.PathUnescape(c.Params("password"))
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}

	// Check if user exists
	res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=\"?\" AND password=\"?\" LIMIT 1", name, password)
	if err != nil {
		return err
	}
	defer res.Close()
	var count int
	err = res.Scan(&count)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	if count == 0 {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    "Invalid username or password",
		})
	}

	res, err = e.db.Query("SELECT uid FROM users WHERE name=\"?\" AND password=\"?\" LIMIT 1", name, password)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	defer res.Close()
	var uid string
	err = res.Scan(&uid)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}

	return c.JSON(map[string]interface{}{
		"success": true,
		"data":    uid,
	})
}

func (e *Elemental) newAnonymousUser(c *fiber.Ctx) error {
	count := 1
	var name string
	var err error
	for count != 0 {
		name, err = GenerateRandomStringURLSafe(8)
		if err != nil {
			return c.JSON(map[string]interface{}{
				"success": false,
				"data":    err.Error(),
			})
		}

		// Check if name taken
		res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=\"?\" LIMIT 1", name)
		if err != nil {
			return c.JSON(map[string]interface{}{
				"success": false,
				"data":    err.Error(),
			})
		}
		defer res.Close()
		err = res.Scan(&count)
		if err != nil {
			return c.JSON(map[string]interface{}{
				"success": false,
				"data":    err.Error(),
			})
		}
	}
	return c.JSON(map[string]interface{}{
		"success": true,
		"data":    name,
	})
}
