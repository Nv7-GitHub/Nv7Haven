package single

import (
	"errors"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type listItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UID         string `json:"uid"`
	ID          string `json:"id"`
}

func (s *Single) list(c *fiber.Ctx) error {
	var kinds = map[string]string{
		"date":  "createdOn DESC",
		"az":    "title ASC",
		"likes": "likes DESC",
		"za":    "title DESC",
	}
	query, err := url.PathUnescape(c.Params("query"))
	if err != nil {
		return err
	}
	kind, exists := kinds[c.Params("kind")]
	if !exists {
		return errors.New("invalid kind")
	}
	res, err := s.db.Query("SELECT title, description, uid, id FROM single WHERE title LIKE ? ORDER BY "+kind, query+"%")
	if err != nil {
		return err
	}
	defer res.Close()
	list := make([]listItem, 0)
	for res.Next() && len(list) < 11 {
		item := listItem{}
		err = res.Scan(&item.Title, &item.Description, &item.UID, &item.ID)
		if err != nil {
			return err
		}
		list = append(list, item)
	}
	return c.JSON(list)
}
