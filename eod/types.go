package eod

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type serverDataType int
type pollType int

const (
	playChannel   = 0
	votingChannel = 1
	newsChannel   = 2
	voteCount     = 3

	pollCombo      = 0
	pollCategorize = 1
	pollSign       = 2
	pollImage      = 3
)

type empty struct{}

type serverData struct {
	playChannels  map[string]empty // channelID
	votingChannel string
	newsChannel   string
	voteCount     int
	combCache     map[string]comb             // map[userID]comb
	invCache      map[string]map[string]empty // map[userID]map[elementName]empty
	elemCache     map[string]element          //map[elementName]element
	polls         map[string]poll             // map[messageid]poll
}

type comb struct {
	elem1 string
	elem2 string
	elem3 string
}

type element struct {
	Name       string
	Category   string
	Image      string
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
}

type poll struct {
	Channel string
	Message string
	Guild   string
	Kind    pollType
	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Data    map[string]interface{}

	Upvotes   int
	Downvotes int
}

type msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type rsp interface {
	Error(err error) bool
	ErrorMessage(msg string)
	Message(msg string) string
	Embed(emb *discordgo.MessageEmbed) string
	Resp(msg string)
}
