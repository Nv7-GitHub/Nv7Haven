package polls

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

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
	lasted := fmt.Sprintf(db.Config.LangProperty("Lasted"), time.Since(p.CreatedOn.Time).Round(time.Second).String()) + " • "
	fmt.Println(lasted)

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
			b.dg.ChannelMessageSend(db.Config.NewsChannel, fmt.Sprintf(db.Config.LangProperty("AddCatNews"), name.Name, p.PollCategorizeData.Category, lasted, p.Suggestor)+controversialTxt)
		} else {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, fmt.Sprintf(db.Config.LangProperty("AddCatMultNews"), len(els), p.PollCategorizeData.Category, lasted, p.Suggestor)+controversialTxt)
		}
	case types.PollUnCategorize:
		els := p.PollCategorizeData.Elems
		for _, val := range els {
			b.UnCategorize(val, p.PollCategorizeData.Category, p.Guild)
		}
		if len(els) == 1 {
			name, _ := db.GetElement(els[0])
			b.dg.ChannelMessageSend(db.Config.NewsChannel, fmt.Sprintf(db.Config.LangProperty("RmCatNews"), name.Name, p.PollCategorizeData.Category, lasted, p.Suggestor)+controversialTxt)
		} else {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, fmt.Sprintf(db.Config.LangProperty("RmCatMultNews"), len(els), p.PollCategorizeData.Category, lasted, p.Suggestor)+controversialTxt)
		}
	case types.PollCatImage:
		b.catImage(p.Guild, p.PollCatImageData.Category, p.PollCatImageData.NewImage, p.Suggestor, p.PollCatImageData.Changed, controversialTxt, lasted, true)
	case types.PollColor:
		b.color(p.Guild, p.PollColorData.Element, p.PollColorData.Color, p.Suggestor, controversialTxt, lasted, true)
	case types.PollCatColor:
		b.catColor(p.Guild, p.PollCatColorData.Category, p.PollCatColorData.Color, p.Suggestor, controversialTxt, lasted, true)
	}
}
