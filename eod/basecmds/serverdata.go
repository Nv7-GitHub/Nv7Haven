package basecmds

import (
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *BaseCmds) SetNewsChannel(channelID string, msg types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=?", msg.GuildID, types.NewsChannel)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}

	if count == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET value1=? WHERE guild=? AND type=?", channelID, msg.GuildID, types.NewsChannel)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.NewsChannel, channelID, 0)
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	dat.NewsChannel = channelID
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully updated news channel!")
}

func (b *BaseCmds) SetVotingChannel(channelID string, msg types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=?", msg.GuildID, types.VotingChannel)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}

	if count == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET value1=? WHERE guild=? AND type=?", channelID, msg.GuildID, types.VotingChannel)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.VotingChannel, channelID, 0)
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	dat.VotingChannel = channelID
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully updated voting channel!")
}

func (b *BaseCmds) SetVoteCount(count int, msg types.Msg, rsp types.Rsp) {
	if count < 0 {
		count *= -1
	}
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=?", msg.GuildID, types.VoteCount)
	var cnt int
	err := row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	if cnt == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET intval=? WHERE guild=? AND type=?", count, msg.GuildID, types.VoteCount)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.VoteCount, "", count)
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	dat.VoteCount = count
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully updated vote count!")
}

func (b *BaseCmds) SetPollCount(count int, msg types.Msg, rsp types.Rsp) {
	if count < 0 {
		count *= -1
	}
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=?", msg.GuildID, types.PollCount)
	var cnt int
	err := row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	if cnt == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET intval=? WHERE guild=? AND type=?", count, msg.GuildID, types.PollCount)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.PollCount, "", count)
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	dat.PollCount = count
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully updated poll count!")
}

func (b *BaseCmds) SetPlayChannel(channelID string, isPlayChannel bool, msg types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", msg.GuildID, types.PlayChannel, channelID)
	var cnt int
	err := row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	if cnt == 1 && !isPlayChannel {
		_, err = b.db.Exec("DELETE FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", msg.GuildID, types.PlayChannel, channelID)
		if rsp.Error(err) {
			return
		}

		b.lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		b.lock.RUnlock()
		if !exists {
			dat = types.NewServerDat()
		}
		delete(dat.PlayChannels, channelID)
		b.lock.Lock()
		b.dat[msg.GuildID] = dat
		b.lock.Unlock()

		rsp.Message("Succesfully marked channel as not a play channel.")
		return
	}

	if !isPlayChannel {
		rsp.ErrorMessage("Channel isn't play channel!")
		return
	}

	_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.PlayChannel, channelID, 0)
	if rsp.Error(err) {
		return
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	if dat.PlayChannels == nil {
		dat.PlayChannels = make(map[string]types.Empty)
	}
	dat.PlayChannels[channelID] = types.Empty{}
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully marked channel as play channel!")
}

func (b *BaseCmds) SetModRole(roleID string, msg types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=?", msg.GuildID, types.ModRole)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}

	if count == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET value1=? WHERE guild=? AND type=?", roleID, msg.GuildID, types.ModRole)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", msg.GuildID, types.ModRole, roleID, 0)
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	dat.ModRole = roleID
	b.lock.Lock()
	b.dat[msg.GuildID] = dat
	b.lock.Unlock()

	rsp.Message("Succesfully updated mod role!")
}

func (b *BaseCmds) SetUserColor(color string, removeColor bool, m types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", m.GuildID, types.UserColor, m.Author.ID)
	var cnt int
	err := row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	// Remove color
	if cnt == 1 && removeColor {
		_, err = b.db.Exec("DELETE FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", m.GuildID, types.UserColor, m.Author.ID)
		if rsp.Error(err) {
			return
		}

		b.lock.RLock()
		dat, exists := b.dat[m.GuildID]
		b.lock.RUnlock()
		if !exists {
			rsp.ErrorMessage("Guild not set up yet!")
			return
		}
		delete(dat.UserColors, m.Author.ID)
		b.lock.Lock()
		b.dat[m.GuildID] = dat
		b.lock.Unlock()

		rsp.Message("Successfully reset color!")
		return
	}

	if removeColor {
		rsp.ErrorMessage("You don't have a color!")
		return
	}

	// Parse
	if len(color) > 0 && color[0] == '#' {
		color = color[1:]
	}
	if len(color) != 6 {
		rsp.ErrorMessage("A hex color must be 6 characters long!")
		return
	}
	col, err := strconv.ParseInt(color, 16, 64)
	if rsp.Error(err) {
		return
	}

	// Update
	if cnt == 1 {
		_, err = b.db.Exec("UPDATE eod_serverdata SET intval=? WHERE guild=? AND type=? AND value1=?", int(col), m.GuildID, types.UserColor, m.Author.ID)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", m.GuildID, types.UserColor, m.Author.ID, int(col))
		if rsp.Error(err) {
			return
		}
	}

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		dat = types.NewServerDat()
	}
	if dat.PlayChannels == nil {
		dat.PlayChannels = make(map[string]types.Empty)
	}
	dat.UserColors[m.Author.ID] = int(col)
	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()

	if cnt == 0 {
		rsp.Message("Successfully set color!")
	} else {
		rsp.Message("Successfully updated color!")
	}
}
