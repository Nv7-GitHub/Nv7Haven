package nv7haven

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jdkato/prose/v2"
)

func (d *Nv7Haven) calcHella(c *fiber.Ctx) error {

	input, err := url.PathUnescape(c.Params("input"))
	if err != nil {
		return err
	}
	doc, _ := prose.NewDocument(input)

	done := []string{}

	// Iterate over the doc's tokens:
	for _, tok := range doc.Tokens() {
		if tok.Tag == "JJ" || tok.Tag == "JJR" || tok.Tag == "JJS" {
			if !(isIn(tok.Tag, done)) {
				done = append(done, tok.Tag)
				input = strings.Replace(input, tok.Text, "hella-"+tok.Text, -1)
			}
		}
	}

	return c.SendString(input)
}

func isIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
