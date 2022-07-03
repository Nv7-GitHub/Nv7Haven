package basecmds

import (
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *BaseCmds) SetNewsChannel(channelID string, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.NewsChannel = channelID
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("NewsChannel", nil))
}

func (b *BaseCmds) SetVotingChannel(channelID string, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.VotingChannel = channelID
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("VotingChannel", nil))
}

func (b *BaseCmds) SetVoteCount(count int, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.VoteCount = count
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("VoteCount", nil))
}

func (b *BaseCmds) SetPollCount(count int, msg types.Msg, rsp types.Rsp) {
	if count < 0 {
		count *= -1
	}
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.PollCount = count
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("PollCount", nil))
}

func (b *BaseCmds) SetPlayChannel(channelID string, isPlayChannel bool, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	if isPlayChannel {
		db.Config.Lock()
		db.Config.PlayChannels[channelID] = types.Empty{}
		db.Config.Unlock()
	} else {
		db.Config.Lock()
		delete(db.Config.PlayChannels, channelID)
		db.Config.Unlock()
	}

	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	if isPlayChannel {
		rsp.Message(db.Config.LangProperty("PlayChannelNew", nil))
	} else {
		rsp.Message(db.Config.LangProperty("PlayChannelRemove", nil))
	}
}

func (b *BaseCmds) SetModRole(roleID string, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.ModRole = roleID
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("ModRole", nil))
}

func (b *BaseCmds) SetUserColor(color string, removeColor bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	// Remove color
	if removeColor {
		db.Config.Lock()
		delete(db.Config.UserColors, m.Author.ID)
		db.Config.Unlock()
		err := db.SaveConfig()
		if rsp.Error(err) {
			return
		}
		rsp.Message(db.Config.LangProperty("UserColorReset", nil))
		return
	}

	// Parse
	if len(color) > 0 && color[0] == '#' {
		color = color[1:]
	}
	if len(color) != 6 {
		rsp.ErrorMessage(db.Config.LangProperty("HexMustBe6", nil))
		return
	}
	col, err := strconv.ParseInt(color, 16, 64)
	if rsp.Error(err) {
		return
	}

	// Update
	db.Config.Lock()
	db.Config.UserColors[m.Author.ID] = int(col)
	db.Config.Unlock()

	err = db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("UserColor", nil))
}

func (b *BaseCmds) SetLanguage(lang string, msg types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(msg.GuildID)
	if !res.Exists {
		var err error
		db, err = b.NewDB(msg.GuildID)
		if rsp.Error(err) {
			return
		}
	}
	db.Config.LanguageFile = lang
	err := db.SaveConfig()
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("Language", nil))
}
