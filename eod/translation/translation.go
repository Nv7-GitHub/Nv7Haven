package translation

import (
	"bytes"
	"embed"
	"encoding/json"
	"sort"
	"text/template"
)

//go:embed languages/*.json
var langData embed.FS
var langFiles = make(map[string]translation)

const DefaultLang = "en_us"

type translation map[string]*template.Template

func mustExecute(tmpl *template.Template, params interface{}) string {
	out := bytes.NewBuffer(nil)
	err := tmpl.Execute(out, params)
	if err != nil {
		panic(err)
	}
	return out.String()
}

func init() {
	files, err := langData.ReadDir("languages")
	if err != nil {
		panic(err)
	}

	var langDat map[string]string
	for _, file := range files {
		f, err := langData.Open("languages/" + file.Name())
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(f)
		err = dec.Decode(&langDat)
		if err != nil {
			panic(err)
		}

		lang := make(translation, len(langDat))
		for k, v := range langDat {
			lang[k] = template.Must(template.New(k).Parse(v))
		}
		langFiles[file.Name()[:len(file.Name())-5]] = lang // Remove .json

		f.Close()
		langDat = nil
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
			Name: mustExecute(f["Name"], nil),
			Lang: lang,
		}
		i++
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Name < langs[j].Name
	})
	return langs
}

func LangProperty(lang, property string, params interface{}) string {
	v, exists := langFiles[lang][property]
	if !exists {
		return mustExecute(langFiles[DefaultLang][property], params)
	}
	return mustExecute(v, params)
}
