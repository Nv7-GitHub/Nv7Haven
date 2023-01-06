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
