package nv7haven

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type idea struct {
	ID        int
	CreatedOn int
	Yes       int
	No        int
	Title     string
	HasVoted  bool
}

type empty struct{}

func (n *Nv7Haven) getIdeas(c *fiber.Ctx) error {
	sort := "votes DESC"
	sortVal := c.Params("sort")
	if sortVal == "new" {
		sort = "createdOn DESC"
	}

	res, err := n.sql.Query("SELECT id, createdOn, yes, no, text, voted FROM ideas WHERE 1 ORDER BY " + sort)
	if err != nil {
		return err
	}
	defer res.Close()

	ip := c.IPs()[0]
	var voted string
	var votes map[string]empty

	out := make([]idea, 0)
	for res.Next() {
		val := idea{}
		res.Scan(&val.ID, &val.CreatedOn, &val.Yes, &val.No, &val.Title, &voted)
		votes = make(map[string]empty)
		err = json.Unmarshal([]byte(voted), &votes)
		if err != nil {
			return err
		}
		_, val.HasVoted = votes[ip]

		out = append(out, val)
	}
	return c.JSON(out)
}

func (n *Nv7Haven) newIdea(c *fiber.Ctx) error {
	text, err := url.PathUnescape(c.Params("title"))
	if err != nil {
		return err
	}
	_, err = n.sql.Exec("INSERT INTO ideas VALUES (?, ?, ?, ?, ?, ?, ?)", rand.Intn(1000000), time.Now().Unix(), 0, 0, 0, "{}", text)
	if err != nil {
		return err
	}
	return nil
}

func (n *Nv7Haven) updateIdea(c *fiber.Ctx) error {
	id := c.Params("id")
	vote := true
	if c.Params("vote") == "0" {
		vote = false
	}

	res := n.sql.QueryRow("SELECT yes, no, voted FROM ideas WHERE id=?", id)
	var yes int
	var no int
	var voted string
	err := res.Scan(&yes, &no, &voted)
	if err != nil {
		return err
	}

	ip := c.IPs()[0]
	var votes map[string]empty
	err = json.Unmarshal([]byte(voted), &votes)
	if err != nil {
		return err
	}
	_, hasVoted := votes[ip]
	if hasVoted {
		return c.SendString("You already voted!")
	}
	votes[ip] = empty{}
	votedDat, err := json.Marshal(votes)
	if err != nil {
		return err
	}
	voted = string(votedDat)

	if vote {
		yes++
	} else {
		no++
	}
	_, err = n.sql.Exec("UPDATE ideas SET yes=?, no=?, votes=?, voted=? WHERE id=?", yes, no, yes+no, voted, id)
	if err != nil {
		return err
	}
	return c.SendString(strconv.Itoa(yes) + "\n" + strconv.Itoa(no))
}
