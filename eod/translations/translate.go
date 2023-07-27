package translations

import (
	"github.com/Nv7-Github/sevcord/v2"
        "github.com/Nv7-Github/Nv7Haven/eod/types"
        "bytes"
	"log"
)

func (t *Translations) SetTranslate(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	lang := opts[0].(string)
	_, err := t.db.Exec("UPDATE config SET lang=$1 WHERE user=$2", lang, c.Author().User.ID)
	if err != nil {
		t.base.Error(c, err)
		return
	}

	c.Respond(sevcord.NewMessage("Successfully set language to " + lang + "! ✅"))
}

func Translate(phrase string, var1 interface{}, var2 interface{}, var3 interface{}) string {
	lang := `SELECT FROM config WHERE "user"=$2`
        LanguageTable.EnglishLock.RLock()
	index := LanguageTable[lang][phrase]

	t := template.Must(template.New("phrase").Parse(index))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, Variables{var1, var2, var3}); err != nil {
		log.Println("Translation Error:", err)
	}
        LanguageTable.EnglishLock.RUnlock()

	return tpl.String()
}
