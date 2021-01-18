package single

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type empty struct{}

func (s *Single) like(c *fiber.Ctx) error {
	id := c.Params("id")
	uid := c.Params("uid")
	res := s.db.QueryRow("SELECT likes, likedby FROM single WHERE id=? AND uid=?", id, uid)
	var dat string
	var likes int
	err := res.Scan(&likes, &dat)
	if err != nil {
		return err
	}
	var likedby map[string]empty
	err = json.Unmarshal([]byte(dat), &likedby)
	if err != nil {
		return err
	}
	ip := c.IPs()[0]
	_, exists := likedby[ip]
	if exists {
		return c.SendString("You already liked this!")
	}
	likes++
	likedby[ip] = empty{}
	data, err := json.Marshal(likedby)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE single SET likes=?, likedby=? WHERE id=? AND uid=?", likes, data, id, uid)
	if err != nil {
		return err
	}
	return nil
}
