package nv7haven

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var tfchan map[string]chan string

func (n *Nv7Haven) searchTf(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	query, err := url.PathUnescape(c.Params("query"))
	if err != nil {
		return err
	}

	order, err := url.PathUnescape(c.Params("order"))
	if err != nil {
		return err
	}

	res, err := n.sql.Query("SELECT name FROM tf WHERE createdon>?  AND name LIKE ? ORDER BY "+order+" DESC", time.Now().Add(-24*time.Hour).Unix(), query)
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

func (n *Nv7Haven) newTf(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	_, err := n.sql.Exec("DELETE FROM tf WHERE createdon<?", time.Now().Add(-24*time.Hour).Unix())
	if err != nil {
		return err
	}

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}
	body := string(c.Body())

	_, err = n.sql.Exec("INSERT INTO tf VALUES (?, ?, ?, ?, ?, ?)", name, body, 0, "[]", "[]", time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (n *Nv7Haven) like(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	var likeddat string
	var likes int
	err = n.query("SELECT likes, likedby FROM tf WHERE name=?", []interface{}{name}, &likes, &likeddat)
	if err != nil {
		return err
	}
	var likedby []string
	err = json.Unmarshal([]byte(likeddat), &likedby)
	if err != nil {
		return err
	}
	ip := c.IPs()[0]
	for _, val := range likedby {
		if val == ip {
			return c.SendString("You already liked it!")
		}
	}
	likedby = append(likedby, ip)
	dat, err := json.Marshal(likedby)
	if err != nil {
		return err
	}
	likes++
	_, err = n.sql.Exec("UPDATE tf SET likedby=?, likes=? WHERE name=?", dat, likes, name)
	if err != nil {
		return err
	}
	return nil
}

func (n *Nv7Haven) comment(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	body := string(c.Body())

	var commentDat string
	err = n.query("SELECT comments FROM tf WHERE name=?", []interface{}{name}, &commentDat)
	if err != nil {
		return err
	}

	var comments []string
	err = json.Unmarshal([]byte(commentDat), &comments)
	if err != nil {
		return err
	}
	comments = append(comments, body)

	dat, err := json.Marshal(comments)
	if err != nil {
		return err
	}
	_, err = n.sql.Exec("UPDATE tf SET comments=? WHERE name=?", dat, name)
	if err != nil {
		return err
	}
	_, exists := tfchan[name]
	if exists {
		tfchan[name] <- body
	}
	return nil
}

type post struct {
	Name      string
	Content   string
	Likes     int
	Comments  []string
	CreatedOn int
}

func (n *Nv7Haven) getPost(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	var content string
	var likes int
	var commentdat string
	var createdon int
	err = n.query("SELECT content, likes, comments, createdon FROM tf WHERE name=?", []interface{}{name}, &content, &likes, &commentdat, &createdon)
	if err != nil {
		return err
	}

	var comments []string
	err = json.Unmarshal([]byte(commentdat), &comments)
	if err != nil {
		return err
	}

	return c.JSON(post{
		Name:      name,
		Content:   content,
		Likes:     likes,
		Comments:  comments,
		CreatedOn: createdon,
	})
}

func (n *Nv7Haven) chatUpdates(c *websocket.Conn) {
	log.Println("chatUpdates")
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		log.Println(err)
	}
	log.Println(name)
	_, exists := tfchan[name]
	if !exists {
		tfchan[name] = make(chan string)
	}
	var mt int
	var val string
	for {
		if mt, _, err = c.ReadMessage(); err != nil {
			log.Println(err)
			break
		}

		//val = <-tfchan[name]
		val = "test"

		if err = c.WriteMessage(mt, []byte(val)); err != nil {
			log.Println(err)
			break
		}
	}
}

func (n *Nv7Haven) postUpdates(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		log.Println(err)
	}
	log.Println(name)
	_, exists := tfchan[name]
	if !exists {
		tfchan[name] = make(chan string)
	}
	val := <-tfchan[name]
	return c.SendString(val)
}
