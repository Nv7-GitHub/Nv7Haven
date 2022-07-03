package eod

import (
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type normalResp struct {
	msg    *discordgo.MessageCreate
	b      *EoD
	typing bool
}

func (n *normalResp) Error(err error) bool {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	if err != nil {
		_, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+err.Error()+" "+types.RedCircle)
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return false
		}
		return true
	}
	return false
}

func (n *normalResp) ErrorMessage(msg string) string {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	m, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+msg+" "+types.RedCircle)
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return ""
	}
	return m.ID
}

func (n *normalResp) Resp(msg string, components ...discordgo.MessageComponent) {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	n.Message(msg, components...)
}

func (n *normalResp) Message(msg string, components ...discordgo.MessageComponent) string {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	var err error
	var m *discordgo.Message

	if len(components) == 0 {
		m, err = n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+msg)
	} else {
		msg := &discordgo.MessageSend{
			Content:    n.msg.Author.Mention() + " " + msg,
			Components: components,
		}
		m, err = n.b.dg.ChannelMessageSendComplex(n.msg.ChannelID, msg)
	}

	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return ""
	}
	return m.ID
}

func (n *normalResp) DM(msg string) {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	channel, err := n.b.dg.UserChannelCreate(n.msg.Author.ID)
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return
	}
	_, err = n.b.dg.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return
	}
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	if emb.Color == 0 {
		color, err := n.b.base.GetColor(n.msg.GuildID, n.msg.Author.ID)
		if err == nil {
			emb.Color = color
		}
	}
	m := &discordgo.MessageSend{
		Embed:      emb,
		Components: components,
	}
	msg, err := n.b.dg.ChannelMessageSendComplex(n.msg.ChannelID, m)
	if err != nil {
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return ""
		}
		return ""
	}
	return msg.ID
}

func (n *normalResp) RawEmbed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	if emb.Color == 0 {
		color, err := n.b.base.GetColor(n.msg.GuildID, n.msg.Author.ID)
		if err == nil {
			emb.Color = color
		}
	}

	msg, err := n.b.dg.ChannelMessageSendComplex(n.msg.ChannelID, &discordgo.MessageSend{
		Embed:      emb,
		Components: components,
	})
	if err != nil {
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return ""
		}
		return ""
	}

	return msg.ID
}

func (n *normalResp) Acknowledge() {
	n.b.dg.ChannelTyping(n.msg.ChannelID)
}

func (n *normalResp) Attachment(text string, files []*discordgo.File) {
	if n.typing {
		n.b.dg.ChannelTyping(n.msg.ChannelID)
	}
	n.b.dg.ChannelMessageSendComplex(n.msg.ChannelID, &discordgo.MessageSend{
		Content: text,
		Files:   files,
	})
}

func (n *normalResp) Modal(modal *discordgo.InteractionResponseData, handler types.ModalHandler) {
	n.ErrorMessage("Modals cannot be used with text commands!")
}

func (b *EoD) newMsgNormal(m *discordgo.MessageCreate) types.Msg {
	return types.Msg{
		Author:    m.Author,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
}

func (b *EoD) newRespNormal(m *discordgo.MessageCreate) types.Rsp {
	return &normalResp{
		msg: m,
		b:   b,
	}
}

type slashResp struct {
	i          *discordgo.InteractionCreate
	b          *EoD
	isFollowup bool
}

func (s *slashResp) Error(err error) bool {
	if err != nil {
		if s.isFollowup {
			_, err := s.b.dg.FollowupMessageCreate(s.i.Interaction, true, &discordgo.WebhookParams{
				Content: err.Error(),
			})
			if err != nil {
				log.SetOutput(logs.DiscordLogs)
				log.Println(err)
				return false
			}
		} else {
			err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: err.Error(),
				},
			})
			if err != nil {
				log.SetOutput(logs.DiscordLogs)
				log.Println(err)
				return false
			}
		}
		return true
	}
	return false
}

func (s *slashResp) ErrorMessage(msg string) string {
	if s.isFollowup {
		m, err := s.b.dg.FollowupMessageCreate(s.i.Interaction, true, &discordgo.WebhookParams{
			Content: msg,
		})
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return ""
		}
		return m.ID
	}

	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
	return ""
}

func (s *slashResp) Resp(msg string, components ...discordgo.MessageComponent) {
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      1 << 6,
			Content:    msg,
			Components: components,
		},
	})
}

func (s *slashResp) Message(msg string, components ...discordgo.MessageComponent) string {
	if s.isFollowup {
		msg, err := s.b.dg.FollowupMessageCreate(s.i.Interaction, true, &discordgo.WebhookParams{
			Content:    msg,
			Components: components,
		})
		if err != nil {
			if err != nil {
				log.SetOutput(logs.DiscordLogs)
				log.Println(err)
				return ""
			}
			return ""
		}
		return msg.ID
	}
	err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    msg,
			Components: components,
		},
	})
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return ""
	}
	return ""
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	if emb.Color == 0 {
		color, err := bot.base.GetColor(s.i.GuildID, s.i.Member.User.ID)
		if err == nil {
			emb.Color = color
		}
	}
	if s.isFollowup {
		msg, err := s.b.dg.FollowupMessageCreate(s.i.Interaction, true, &discordgo.WebhookParams{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		})
		if err != nil {
			if err != nil {
				log.SetOutput(logs.DiscordLogs)
				log.Println(err)
				return ""
			}
			return ""
		}
		return msg.ID
	}
	err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		},
	})
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
		return ""
	}
	return ""
}

func (s *slashResp) RawEmbed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	return s.Embed(emb, components...)
}

func (s *slashResp) Acknowledge() {
	s.isFollowup = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func (s *slashResp) DM(msg string) {
	channel, err := s.b.dg.UserChannelCreate(s.i.Member.User.ID)
	if err != nil {
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return
		}
	}
	_, err = s.b.dg.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		if err != nil {
			log.SetOutput(logs.DiscordLogs)
			log.Println(err)
			return
		}
	}
}

func (s *slashResp) Attachment(text string, files []*discordgo.File) {
	if s.isFollowup {
		s.b.dg.FollowupMessageCreate(s.i.Interaction, true, &discordgo.WebhookParams{
			Content: text,
			Files:   files,
		})
	}
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
			Files:   files,
		},
	})
}

func (s *slashResp) Modal(modal *discordgo.InteractionResponseData, handler types.ModalHandler) {
	modal.CustomID = s.i.Interaction.ID
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
	dat, res := s.b.GetData(s.i.GuildID)
	if res.Exists {
		dat.AddModal(s.i.ID, handler)
	}
}

func (b *EoD) newMsgSlash(i *discordgo.InteractionCreate) types.Msg {
	return types.Msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
}

func (b *EoD) newRespSlash(i *discordgo.InteractionCreate) types.Rsp {
	return &slashResp{
		i: i,
		b: b,
	}
}
