package discord

import "github.com/bwmarrin/discordgo"

type normalResp struct {
	msg *discordgo.MessageCreate
	b   *Bot
}

func (n *normalResp) Error(err error) bool {
	if err != nil {
		n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

func (n *normalResp) ErrorMessage(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+msg)
}

func (n *normalResp) Resp(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, msg)
}

func (n *normalResp) Message(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, msg)
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed) {
	n.b.dg.ChannelMessageSendEmbed(n.msg.ChannelID, emb)
}

func (b *Bot) newMsgNormal(m *discordgo.MessageCreate) msg {
	return msg{
		Author:    m.Author,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
}

func (b *Bot) newRespNormal(m *discordgo.MessageCreate) rsp {
	return &normalResp{
		msg: m,
		b:   b,
	}
}

type slashResp struct {
	i *discordgo.InteractionCreate
	b *Bot
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
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: "Error: " + msg,
		},
	})
}

func (s *slashResp) Resp(msg string) {
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
}

func (s *slashResp) Message(msg string) {
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: msg,
		},
	})
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed) {
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{emb},
		},
	})
}

func (b *Bot) newMsgSlash(i *discordgo.InteractionCreate) msg {
	return msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
}

func (b *Bot) newRespSlash(i *discordgo.InteractionCreate) rsp {
	return &slashResp{
		i: i,
		b: b,
	}
}
