package translation

import (
	"embed"
	"encoding/json"
	"sort"
)

//go:embed languages/*.json
var langData embed.FS
var langFiles = make(map[string]translation)

const DefaultLang = "en"

type translation map[string]string

func init() {
	files, err := langData.ReadDir("languages")
	if err != nil {
		panic(err)
	}

	var lang translation
	for _, file := range files {
		f, err := langData.Open("languages/" + file.Name())
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(f)
		err = dec.Decode(&lang)
		if err != nil {
			panic(err)
		}

		langFiles[file.Name()[:len(file.Name())-5]] = lang // Remove .json

		f.Close()
	}
}

type LangFileListItem struct {
	Name string
	Lang string
}

func LangFileList() []LangFileListItem {
	langs := make([]LangFileListItem, len(langFiles))
	i := 0
	for lang, f := range langFiles {
		langs[i] = LangFileListItem{
			Name: f["Name"],
			Lang: lang,
		}
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Name < langs[j].Name
	})
	return langs
}

func LangProperty(lang, property string) string {
	v, exists := langFiles[lang][property]
	if !exists {
		return langFiles[DefaultLang][property]
	}
	return v
}
