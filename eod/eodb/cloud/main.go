package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/psykhi/wordclouds"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const minCnt = 0

var punctuation = []string{"(", ")", ".", ",", ";", "+", "*", ":", "!", "?", "\"", "'", "`", "~", "^", "&", "|", "\\", "/", "=", "<", ">"}

func main() {
	// Load
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading DB...")
	start := time.Now()
	dat, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	// Get counts
	start = time.Now()
	wordCnts := make(map[string]int)
	for _, db := range dat.DB {
		for _, el := range db.Elements {
			words := strings.Split(el.Name, " ")
			for i, word := range words {
				words[i] = strings.ToLower(word)
				for _, char := range punctuation {
					words[i] = strings.ReplaceAll(words[i], char, "")
				}
			}

			// Add
			for _, word := range words {
				wordCnts[word]++
			}
		}
	}

	// Remove below min
	for word, cnt := range wordCnts {
		if cnt < minCnt {
			delete(wordCnts, word)
		}
	}

	fmt.Println("Got counts in", time.Since(start))

	// Render
	fmt.Println("Rendering...")
	start = time.Now()
	w := wordclouds.NewWordcloud(wordCnts,
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.FontFile("./Roboto.ttf"),
		wordclouds.Colors([]color.Color{
			color.RGBA{0x1b, 0x1b, 0x1b, 0xff},
			color.RGBA{0x48, 0x48, 0x4B, 0xff},
			color.RGBA{0x59, 0x3a, 0xee, 0xff},
			color.RGBA{0x65, 0xCD, 0xFA, 0xff},
			color.RGBA{0x70, 0xD6, 0xBF, 0xff},
		}))
	img := w.Draw()
	fmt.Println("Rendered in", time.Since(start))

	// Save
	start = time.Now()

	out, err := os.Create("cloud.png")
	handle(err)
	defer out.Close()

	err = png.Encode(out, img)
	handle(err)

	outJson, err := os.Create("cloud.json")
	handle(err)
	defer out.Close()

	// Sort
	type outJsonVal struct {
		Word  string `json:"word"`
		Count int    `json:"count"`
	}
	outVals := make([]outJsonVal, len(wordCnts))
	i := 0
	for k, v := range wordCnts {
		outVals[i] = outJsonVal{Word: k, Count: v}
		i++
	}
	sort.Slice(outVals, func(i, j int) bool {
		return outVals[i].Count > outVals[j].Count
	})

	enc := json.NewEncoder(outJson)
	enc.SetIndent("", "\t")
	handle(enc.Encode(outVals))
	fmt.Println("Saved in", time.Since(start))
}
