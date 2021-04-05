package eod

import (
	"github.com/bwmarrin/discordgo"
)

const redCircle = "🔴"

type normalResp struct {
	msg *discordgo.MessageCreate
	b   *EoD
}

func (n *normalResp) Error(err error) bool {
	if err != nil {
		n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

func (n *normalResp) ErrorMessage(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+msg+" "+redCircle)
}

func (n *normalResp) Resp(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, n.msg.Author.Mention()+" "+" "+msg)
}

func (n *normalResp) Message(msg string) string {
	m, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, msg)
	if err != nil {
	  return ""
	}
	return m.ID
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed) string {
	msg, err := n.b.dg.ChannelMessageSendEmbed(n.msg.ChannelID, emb)
	if err != nil {
	  return ""
	}
	return msg.ID
}

func (n *normalResp) EmbedFollowup(emb *discordgo.MessageEmbed) string {
	return n.Embed(emb)
}

func (b *EoD) newMsgNormal(m *discordgo.MessageCreate) msg {
	return msg{
		Author:    m.Author,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
}

func (b *EoD) newRespNormal(m *discordgo.MessageCreate) rsp {
	return &normalResp{
		msg: m,
		b:   b,
	}
}

type slashResp struct {
	i          *discordgo.InteractionCreate
	b          *EoD
	hasReplied bool
}

func (s *slashResp) Error(err error) bool {
	if err != nil {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Flags:   1 << 6,
				Content: "Error: " + err.Error(),
			},
		})
		return true
	}
	return false
}

func (s *slashResp) ErrorMessage(msg string) {
	s.hasReplied = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: "Error: " + msg,
		},
	})
}

func (s *slashResp) Resp(msg string) {
	s.hasReplied = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
}

func (s *slashResp) Message(msg string) string {
	if !s.hasReplied {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: msg,
			},
		})
		s.hasReplied = true
		return ""
	}
	m, _ := s.b.dg.ChannelMessageSend(s.i.ChannelID, msg)
	return m.ID
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed) string {
	if !s.hasReplied {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Embeds: []*discordgo.MessageEmbed{emb},
			},
		})
		s.hasReplied = true
		return ""
	}
	m, _ := s.b.dg.ChannelMessageSendEmbed(s.i.ChannelID, emb)
	return m.ID
}

func (s *slashResp) EmbedFollowup(emb *discordgo.MessageEmbed) string {
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	msg, _ := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{emb},
	})
	return msg.ID
}

func (b *EoD) newMsgSlash(i *discordgo.InteractionCreate) msg {
	return msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
}

func (b *EoD) newRespSlash(i *discordgo.InteractionCreate) rsp {
	return &slashResp{
		i: i,
		b: b,
	}
}
