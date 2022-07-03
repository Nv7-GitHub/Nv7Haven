package polls

import (
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) getLasted(db *eodb.DB, p types.Poll) string {
	lasted := ""
	if p.CreatedOn != nil {
		lasted = db.Config.LangProperty("Lasted", time.Since(p.CreatedOn.Time).Round(time.Second).String()) + " • "
	}
	return lasted
}

func (b *Polls) handlePollSuccess(p types.Poll) {
	db, res := b.GetDB(p.Guild)
	if !res.Exists {
		return
	}

	controversial := db.Config.VoteCount != 0 && float32(p.Downvotes)/float32(db.Config.VoteCount) >= 0.3
	controversialTxt := ""
	if controversial {
		controversialTxt = " 🌩️"
	}

	lasted := b.getLasted(db, p)

	switch p.Kind {
	case types.PollCombo:
		b.elemCreate(p.PollComboData.Result, p.PollComboData.Elems, p.Suggestor, controversialTxt, lasted, p.Guild)
	case types.PollSign:
		b.mark(p.Guild, p.PollSignData.Elem, p.PollSignData.NewNote, p.Suggestor, controversialTxt, lasted, true)
	case types.PollImage:
		b.image(p.Guild, p.PollImageData.Elem, p.PollImageData.NewImage, p.Suggestor, p.PollImageData.Changed, controversialTxt, lasted, true)
	case types.PollCategorize:
		els := p.PollCategorizeData.Elems
		for _, val := range els {
			b.Categorize(val, p.PollCategorizeData.Category, p.Guild)
		}
		if len(els) == 1 {
			name, _ := db.GetElement(els[0])
			b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("AddCatNews", map[string]any{
				"Element":    name.Name,
				"Category":   p.PollCategorizeData.Category,
				"LastedText": lasted,
				"Creator":    p.Suggestor,
			})+controversialTxt)
		} else {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("AddCatMultNews", map[string]any{
				"Elements":   len(els),
				"Category":   p.PollCategorizeData.Category,
				"LastedText": lasted,
				"Creator":    p.Suggestor,
			})+controversialTxt)
		}
	case types.PollUnCategorize:
		els := p.PollCategorizeData.Elems
		b.UnCategorize(els, p.PollCategorizeData.Category, p.Guild)
		if len(els) == 1 {
			name, _ := db.GetElement(els[0])
			b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("RmCatNews", map[string]any{
				"Element":    name.Name,
				"Category":   p.PollCategorizeData.Category,
				"LastedText": lasted,
				"Creator":    p.Suggestor,
			})+controversialTxt)
		} else {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("RmCatMultNews", map[string]any{
				"Elements":   len(els),
				"Category":   p.PollCategorizeData.Category,
				"LastedText": lasted,
				"Creator":    p.Suggestor,
			})+controversialTxt)
		}
	case types.PollCatImage:
		b.catImage(p.Guild, p.PollCatImageData.Category, p.PollCatImageData.NewImage, p.Suggestor, p.PollCatImageData.Changed, controversialTxt, lasted, true)
	case types.PollColor:
		b.color(p.Guild, p.PollColorData.Element, p.PollColorData.Color, p.Suggestor, controversialTxt, lasted, true)
	case types.PollCatColor:
		b.catColor(p.Guild, p.PollCatColorData.Category, p.PollCatColorData.Color, p.Suggestor, controversialTxt, lasted, true)
	case types.PollCatSign:
		b.catSign(p.Guild, p.PollCatSignData.CatName, p.PollCatSignData.NewNote, p.Suggestor, controversialTxt, lasted, true)
	case types.PollDeleteVCat:
		b.deleteVCat(p.Guild, p.PollVCatDeleteData.Category, p.Suggestor, controversialTxt, lasted, true)
	}
}
