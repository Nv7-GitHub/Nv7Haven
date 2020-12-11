package elemental

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) createUser(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	email, err := url.PathUnescape(c.Params("email"))
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
	user, err := e.auth.SignUpWithEmailAndPassword(email, password)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	err = e.checkUser(user.OtherData.LocalID)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	return c.JSON(map[string]interface{}{
		"success": true,
		"data":    user.OtherData.LocalID,
	})
}

func (e *Elemental) loginUser(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	email, err := url.PathUnescape(c.Params("email"))
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
	user, err := e.auth.SignInWithEmailAndPassword(email, password)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	err = e.checkUser(user.OtherData.LocalID)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	return c.JSON(map[string]interface{}{
		"success": true,
		"data":    user.OtherData.LocalID,
	})
}

func (e *Elemental) checkUser(uid string) error {
	data, err := e.db.Get("users/" + uid)
	if err != nil {
		return err
	}
	if string(data) == "null" {
		e.db.SetData("users/"+uid+"/found", []string{"Air", "Earth", "Fire", "Water"})
	}
	return nil
}

func (e *Elemental) resetPassword(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	email, err := url.PathUnescape(c.Params("email"))
	if err != nil {
		return err
	}
	err = e.auth.ResetPassword(email)
	if err != nil {
		return err
	}
	return nil
}
