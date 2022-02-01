package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/finnbear/moderation"
)

func (b *Polls) MarkCmd(elem string, mark string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), elem))
		return
	}

	inv := db.GetInv(m.Author.ID)
	exists := inv.Contains(el.ID)
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DontHave"), el.Name))
		return
	}
	if len(mark) >= 2400 {
		rsp.ErrorMessage(db.Config.LangProperty("MaxMarkLength"))
		return
	}
	if len(mark) == 0 {
		mark = db.Config.LangProperty("DefaultMark")
	}
	if moderation.IsInappropriate(mark) && db.Config.SwearFilter {
		rsp.ErrorMessage(db.Config.LangProperty("NoInappropriateSuggest"))
		return
	}

	if el.Creator == m.Author.ID {
		id, res := db.GetIDByName(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		b.mark(m.GuildID, id, mark, m.Author.ID, "", false)
		rsp.Message(fmt.Sprintf(db.Config.LangProperty("MarkChanged"), el.Name))
		return
	}

	err := b.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollSign,
		Suggestor: m.Author.ID,

		PollSignData: &types.PollSignData{
			Elem:    el.ID,
			NewNote: mark,
			OldNote: el.Comment,
		},
	})
	if rsp.Error(err) {
		return
	}
	id := rsp.Message(fmt.Sprintf(db.Config.LangProperty("MarkSuggested"), el.Name))

	data, _ := b.GetData(m.GuildID)
	data.SetMsgElem(id, el.ID)
}

func (b *Polls) ImageCmd(elem string, image string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), elem))
		return
	}

	inv := db.GetInv(m.Author.ID)
	exists := inv.Contains(el.ID)
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DontHave"), el.Name))
		return
	}

	changed := el.Image != ""

	if el.Creator == m.Author.ID {
		id, res := db.GetIDByName(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		b.image(m.GuildID, id, image, m.Author.ID, changed, "", false)
		if !changed {
			rsp.Message(fmt.Sprintf(db.Config.LangProperty("ImageAdded"), el.Name))
		} else {
			rsp.Message(fmt.Sprintf(db.Config.LangProperty("ImageChanged"), el.Name))
		}
		return
	}

	err := b.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollImage,
		Suggestor: m.Author.ID,

		PollImageData: &types.PollImageData{
			Elem:     el.ID,
			NewImage: image,
			OldImage: el.Image,
			Changed:  changed,
		},
	})
	if rsp.Error(err) {
		return
	}
	id := rsp.Message(fmt.Sprintf(db.Config.LangProperty("ImageSuggested"), el.Name))
	data, _ := b.GetData(m.GuildID)
	data.SetMsgElem(id, el.ID)
}

func (b *Polls) ColorCmd(elem string, color int, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), elem))
		return
	}

	inv := db.GetInv(m.Author.ID)
	exists := inv.Contains(el.ID)
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DontHave"), el.Name))
		return
	}

	if el.Creator == m.Author.ID {
		id, res := db.GetIDByName(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		b.color(m.GuildID, id, color, m.Author.ID, "", false)
		rsp.Message(fmt.Sprintf(db.Config.LangProperty("ElemColorChanged"), el.Name))
		return
	}

	err := b.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollColor,
		Suggestor: m.Author.ID,

		PollColorData: &types.PollColorData{
			Element:  el.ID,
			Color:    color,
			OldColor: el.Color,
		},
	})
	if rsp.Error(err) {
		return
	}
	id := rsp.Message(fmt.Sprintf(db.Config.LangProperty("ElemColorSuggested"), el.Name))
	data, _ := b.GetData(m.GuildID)
	data.SetMsgElem(id, el.ID)
}

func (b *Polls) CatImgCmd(catName string, url string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	changed := cat.Image != ""

	err := b.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollCatImage,
		Suggestor: m.Author.ID,

		PollCatImageData: &types.PollCatImageData{
			Category: cat.Name,
			NewImage: url,
			OldImage: cat.Image,
			Changed:  changed,
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf(db.Config.LangProperty("CatImageSuggested"), cat.Name))
}

func (b *Polls) CatColorCmd(catName string, color int, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	err := b.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollCatColor,
		Suggestor: m.Author.ID,

		PollCatColorData: &types.PollCatColorData{
			Category: cat.Name,
			Color:    color,
			OldColor: cat.Color,
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf(db.Config.LangProperty("CatColorSuggested"), cat.Name))
}
