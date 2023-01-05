package types

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

type translations struct {
	nameMayNotContain    string
	nameCannotBeEmpty    string
	successfullyUpdated  string
	playChannelReset     string
	youAreNotAuthorized  string
	successUpdateChannel string
}

type Language struct {
	name         string
	translations translations
}

type Variables struct {
	var1 interface{}
	var2 interface{}
	var3 interface{}
}

var LanguageTable = []Language{
	{
		name: "English",
		translations: translations{
			`A name may not contain '{{.var1}}'!`,
			"Name cannot be empty!",
			"Successfully updated {{.var1}}!",
			"**PLAY CHANNELS HAVE BEEN RESET**\nUpdate them below!",
			"You are not authorized to use this!",
			"Successfully updated play channels!",
		},
	},
}

func Translate(phrase string, var1 interface{}, var2 interface{}, var3 interface{}) string {
	lang := `SELECT FROM config WHERE "user"=$2`
	for i := range LanguageTable {
		if LanguageTable[i].name == lang {
			phrase = getAttr(&LanguageTable[i].translations, phrase).String()
			break
		}
	}

	t := template.Must(template.New("phrase").Parse(phrase))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, Variables{var1, var2, var3}); err != nil {
		fmt.Println(err)
	}

	return tpl.String()
}

func getAttr(obj interface{}, fieldName string) reflect.Value {
	pointToStruct := reflect.ValueOf(obj) // addressable
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(fieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField
}
