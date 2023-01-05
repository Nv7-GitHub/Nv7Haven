package translations

import (
	"github.com/Nv7-Github/sevcord/v2"
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
