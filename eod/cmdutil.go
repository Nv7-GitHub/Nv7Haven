package eod

import (
	"log"
	"os"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const redCircle = "🔴"

var discordlogs *os.File

type normalResp struct {
	msg *discordgo.MessageCreate
	b   *EoD
}

func (n *normalResp) Error(err error) bool {
	if err != nil {
		_, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" Error: "+err.Error()+" "+redCircle)
		if err != nil {
			log.SetOutput(discordlogs)
			log.Println(err)
		}
		return true
	}
	return false
}

func (n *normalResp) ErrorMessage(msg string) string {
	m, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+msg+" "+redCircle)
	if err != nil {
		log.SetOutput(discordlogs)
		log.Println(err)
	}
	return m.ID
}

func (n *normalResp) Resp(msg string, components ...discordgo.MessageComponent) {
	n.Message(msg, components...)
}

func (n *normalResp) Message(msg string, components ...discordgo.MessageComponent) string {
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
		log.SetOutput(discordlogs)
		log.Println(err)
		return ""
	}
	return m.ID
}

func (n *normalResp) DM(msg string) {
	channel, err := n.b.dg.UserChannelCreate(n.msg.Author.ID)
	if err != nil {
		log.SetOutput(discordlogs)
		log.Println(err)
	}
	_, err = n.b.dg.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		log.SetOutput(discordlogs)
		log.Println(err)
	}
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	color, err := n.b.getColor(n.msg.GuildID, n.msg.Author.ID)
	if err == nil {
		emb.Color = color
	}
	m := &discordgo.MessageSend{
		Embed:      emb,
		Components: components,
	}
	msg, err := n.b.dg.ChannelMessageSendComplex(n.msg.ChannelID, m)
	if err != nil {
		if err != nil {
			log.SetOutput(discordlogs)
			log.Println(err)
		}
		return ""
	}
	return msg.ID
}

func (n *normalResp) RawEmbed(emb *discordgo.MessageEmbed) string {
	color, err := n.b.getColor(n.msg.GuildID, n.msg.Author.ID)
	if err == nil {
		emb.Color = color
	}

	msg, err := n.b.dg.ChannelMessageSendEmbed(n.msg.ChannelID, emb)
	if err != nil {
		if err != nil {
			log.SetOutput(discordlogs)
			log.Println(err)
		}
		return ""
	}

	return msg.ID
}

func (n *normalResp) Acknowledge() {}

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
			_, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
				Content: "Error: " + err.Error(),
			})
			if err != nil {
				log.SetOutput(discordlogs)
				log.Println(err)
			}
		} else {
			err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Error: " + err.Error(),
				},
			})
			if err != nil {
				log.SetOutput(discordlogs)
				log.Println(err)
			}
		}
		return true
	}
	return false
}

func (s *slashResp) ErrorMessage(msg string) string {
	if s.isFollowup {
		m, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error: " + msg,
		})
		if err != nil {
			log.SetOutput(discordlogs)
			log.Println(err)
		}
		return m.ID
	}

	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: "Error: " + msg,
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
		msg, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Content:    msg,
			Components: components,
		})
		if err != nil {
			if err != nil {
				log.SetOutput(discordlogs)
				log.Println(err)
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
		log.SetOutput(discordlogs)
		log.Println(err)
	}
	return ""
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string {
	color, err := bot.getColor(s.i.GuildID, s.i.Member.User.ID)
	if err == nil {
		emb.Color = color
	}
	if s.isFollowup {
		msg, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		})
		if err != nil {
			if err != nil {
				log.SetOutput(discordlogs)
				log.Println(err)
			}
			return ""
		}
		return msg.ID
	}
	err = s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		},
	})
	if err != nil {
		log.SetOutput(discordlogs)
		log.Println(err)
	}
	return ""
}

func (s *slashResp) RawEmbed(emb *discordgo.MessageEmbed) string {
	return s.Embed(emb)
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
			log.SetOutput(discordlogs)
			log.Println(err)
		}
	}
	_, err = s.b.dg.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		if err != nil {
			log.SetOutput(discordlogs)
			log.Println(err)
		}
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
