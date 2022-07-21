package elements

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

// TODO: Translate

func (b *Elements) EditElementNameCmd(elem string, name string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}
	el.Name = name
	err := db.SaveElement(el)
	if rsp.Error(err) {
		return
	}

	rsp.Message(fmt.Sprintf("Successfully updated element **#%d***!", el.ID))
}

func (b *Elements) EditElementImageCmd(elem string, image string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}
	el.Image = image
	err := db.SaveElement(el)
	if rsp.Error(err) {
		return
	}

	rsp.Message(fmt.Sprintf("Successfully updated element **#%d***!", el.ID))
}

func (b *Elements) EditElementCreatorCmd(elem string, creator string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}
	el.Creator = creator
	err := db.SaveElement(el)
	if rsp.Error(err) {
		return
	}

	rsp.Message(fmt.Sprintf("Successfully updated element **#%d***!", el.ID))
}

func (b *Elements) EditElementCreatedOnCmd(elem string, createdon int64, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}
	el.CreatedOn = types.NewTimeStamp(time.Unix(createdon, 0))
	err := db.SaveElement(el)
	if rsp.Error(err) {
		return
	}

	rsp.Message(fmt.Sprintf("Successfully updated element **#%d***!", el.ID))
}

type editMarkModal struct {
	m    types.Msg
	db   *eodb.DB
	elem int
}

func (m *editMarkModal) Handler(s *discordgo.Session, i *discordgo.InteractionCreate, rsp types.Rsp) {
	el, _ := m.db.GetElement(m.elem)
	el.Comment = i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	err := m.db.SaveElement(el)
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf("Successfully updated element **#%d***!", el.ID))
}

func (b *Elements) EditElementMarkCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}

	rsp.Modal(&discordgo.InteractionResponseData{
		Title: "Mark Element",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "mark",
						Label:       "New Element Mark",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "None",
						Required:    true,
						MinLength:   1,
						MaxLength:   2400,
					},
				},
			},
		},
	}, &editMarkModal{
		m:    m,
		db:   db,
		elem: el.ID,
	})
}
