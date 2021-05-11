package eod

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const redCircle = "🔴"

type normalResp struct {
	msg *discordgo.MessageCreate
	b   *EoD
}

func (n *normalResp) Error(err error) bool {
	if err != nil {
		_, err := n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+err.Error())
		if err != nil {
			log.Println("Failed to send message:", err)
		}
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
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		return ""
	}
	return m.ID
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed) string {
	for _, roleID := range n.msg.Member.Roles {
		role, err := bot.getRole(roleID, n.msg.GuildID)
		if err == nil {
			if role.Color != 0 {
				emb.Color = role.Color
				break
			}
		} else {
			break
		}
	}
	msg, err := n.b.dg.ChannelMessageSendEmbed(n.msg.ChannelID, emb)
	if err != nil {
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		return ""
	}
	return msg.ID
}

func (n *normalResp) Acknowledge() {}

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
	isFollowup bool
}

func (s *slashResp) Error(err error) bool {
	if err != nil {
		if s.isFollowup {
			_, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
				Content: "Error: " + err.Error(),
			})
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		} else {
			err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Flags:   1 << 6,
					Content: "Error: " + err.Error(),
				},
			})
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		}
		return true
	}
	return false
}

func (s *slashResp) ErrorMessage(msg string) {
	if s.isFollowup {
		s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error: " + msg,
		})
	}

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

func (s *slashResp) Message(msg string) string {
	if s.isFollowup {
		msg, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Content: msg,
		})
		if err != nil {
			if err != nil {
				log.Println("Failed to send message:", err)
			}
			return ""
		}
		return msg.ID
	}
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: msg,
		},
	})
	return ""
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed) string {
	for _, roleID := range s.i.Member.Roles {
		role, err := bot.getRole(roleID, s.i.GuildID)
		if err == nil {
			if role.Color != 0 {
				emb.Color = role.Color
				break
			}
		} else {
			break
		}
	}
	if s.isFollowup {
		msg, err := s.b.dg.FollowupMessageCreate(clientID, s.i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{emb},
		})
		if err != nil {
			if err != nil {
				log.Println("Failed to send message:", err)
			}
			return ""
		}
		return msg.ID
	}
	err := s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{emb},
		},
	})
	if err != nil {
		log.Println("Failed to send message:", err)
	}
	return ""
}

func (s *slashResp) Acknowledge() {
	s.isFollowup = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
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
