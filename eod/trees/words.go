package trees

import (
	"image"
	"image/color"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/psykhi/wordclouds"
)

var punctuation = []string{"(", ")", ".", ",", ";", "+", "*", ":", "!", "?", "\"", "'", "`", "~", "^", "&", "|", "\\", "/", "=", "<", ">"}
var colors = []color.Color{
	color.RGBA{0x1b, 0x1b, 0x1b, 0xff},
	color.RGBA{0x48, 0x48, 0x4B, 0xff},
	color.RGBA{0x59, 0x3a, 0xee, 0xff},
	color.RGBA{0x65, 0xCD, 0xFA, 0xff},
	color.RGBA{0x70, 0xD6, 0xBF, 0xff},
}

func getWords(name string) []string {
	words := strings.Split(name, " ")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		for _, char := range punctuation {
			words[i] = strings.ReplaceAll(words[i], char, "")
		}
	}
	return words
}

type WordTree struct {
	wordCnts map[string]int
	added    map[int]types.Empty
	db       *eodb.DB

	CalcTree bool
}

func (w *WordTree) AddElem(elem int, notoplevel ...bool) (bool, string) {
	if len(notoplevel) == 0 {
		w.db.RLock()
		defer w.db.RUnlock()
	}

	_, exists := w.added[elem]
	if exists {
		return true, ""
	}

	el, res := w.db.GetElement(elem, true)
	if !res.Exists {
		return false, res.Message
	}
	words := getWords(el.Name)
	for _, word := range words {
		w.wordCnts[word]++
	}
	w.added[elem] = types.Empty{}

	if w.CalcTree {
		for _, el := range el.Parents {
			w.AddElem(el, true)
		}
	}

	return true, ""
}

func (w *WordTree) Render(width, height int) image.Image {
	cloud := wordclouds.NewWordcloud(w.wordCnts, wordclouds.Width(width), wordclouds.Height(height), wordclouds.FontFile("eod/assets/Roboto.ttf"), wordclouds.Colors(colors))
	return cloud.Draw()
}

func NewWordTree(db *eodb.DB) *WordTree {
	return &WordTree{
		wordCnts: make(map[string]int),
		added:    make(map[int]types.Empty),
		db:       db,
		CalcTree: true,
	}
}
