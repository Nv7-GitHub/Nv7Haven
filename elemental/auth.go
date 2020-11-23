package elemental

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func createUser(c *fiber.Ctx) error {
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
	user, err := auth.SignUpWithEmailAndPassword(email, password)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	err = checkUser(user.OtherData.LocalID)
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

func loginUser(c *fiber.Ctx) error {
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
	user, err := auth.SignInWithEmailAndPassword(email, password)
	if err != nil {
		return c.JSON(map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
	}
	err = checkUser(user.OtherData.LocalID)
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

func checkUser(uid string) error {
	data, err := db.Get("users/" + uid)
	if err != nil {
		return err
	}
	if string(data) == "null" {
		db.SetData("users/"+uid, []string{"Air", "Earth", "Fire", "Water"})
	}
	return nil
}

func resetPassword(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	email, err := url.PathUnescape(c.Params("email"))
	if err != nil {
		return err
	}
	err = auth.ResetPassword(email)
	if err != nil {
		return err
	}
	return nil
}
