package nv7haven

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/jdkato/prose/v2"
)

type section struct {
	Part  string
	Words string
}

var tagMap = map[string]string{
	"(":    "(",
	")":    ")",
	",":    ",",
	":":    ":",
	".":    ".",
	"''":   "”",
	"``":   "“",
	"#":    "#",
	"$":    "$",
	"CC":   "Conjunction",
	"CD":   "Cardinal Number",
	"DT":   "Determiner",
	"EX":   "Existential There",
	"FW":   "Foreign Word",
	"IN":   "Conjunction, Subordinating, or Preposition",
	"JJ":   "Adjective",
	"JJR":  "Adjective, Comparitave",
	"JJS":  "Adjective, Superlative",
	"LS":   "List Item Maker",
	"MD":   "Verb, Modal Auxilary",
	"NN":   "Noun, Singular or Mass",
	"NNP":  "Noun, Proper Singular",
	"NNS":  "Noun, Plural",
	"PDT":  "Predeterminer",
	"POS":  "Possesive Ending",
	"PRP":  "Pronoun, Personal",
	"PRP$": "Pronoun, Possesive",
	"RB":   "Adverb",
	"RBR":  "Adverb, Comparitave",
	"RBS":  "Adverb, Superlative",
	"RP":   "Adverb, Particle",
	"SYM":  "Symbol",
	"TO":   "Infinitival To",
	"UH":   "Interjection",
	"VB":   "Verb, Base Form",
	"VBD":  "Verb, Past Tense",
	"VBG":  "Verb, Gerund or Present Participle",
	"VBN":  "Verb, Past Participle",
	"VBP":  "Verb, Non-3rd Person Singular Present",
	"VBZ":  "Verb, 3rd Person Singular Present",
	"WDT":  "Wh-Determiner",
	"WP":   "Wh-Pronoun, Personal",
	"WP$":  "Wh-Pronoun, Possesive",
	"WRB":  "Wh-Adverb",
}

func (n *Nv7Haven) breakdown(c *fiber.Ctx) error {
	input, err := url.PathUnescape(c.Params("input"))
	if err != nil {
		return err
	}
	doc, _ := prose.NewDocument(input)
	toks := doc.Tokens()
	out := make([]section, len(toks))
	var exists bool
	for i, tok := range toks {
		out[i].Part, exists = tagMap[tok.Tag]
		if !exists {
			out[i].Part = "UNKNOWN"
		}
		out[i].Words = tok.Text
	}
	return c.JSON(out)
}
