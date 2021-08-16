package eod

import (
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) setNewsChannel(channelID string, msg types.Msg, rsp types.Rsp) {
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	dat.NewsChannel = channelID
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully updated news channel!")
}

func (b *EoD) setVotingChannel(channelID string, msg types.Msg, rsp types.Rsp) {
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	dat.VotingChannel = channelID
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully updated voting channel!")
}

func (b *EoD) setVoteCount(count int, msg types.Msg, rsp types.Rsp) {
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	dat.VoteCount = count
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully updated vote count!")
}

func (b *EoD) setPollCount(count int, msg types.Msg, rsp types.Rsp) {
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	dat.PollCount = count
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully updated poll count!")
}

func (b *EoD) setPlayChannel(channelID string, isPlayChannel bool, msg types.Msg, rsp types.Rsp) {
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

		lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		lock.RUnlock()
		if !exists {
			dat = types.NewServerData()
		}
		delete(dat.PlayChannels, channelID)
		lock.Lock()
		b.dat[msg.GuildID] = dat
		lock.Unlock()

		rsp.Resp("Succesfully marked channel as not a play channel.")
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	if dat.PlayChannels == nil {
		dat.PlayChannels = make(map[string]types.Empty)
	}
	dat.PlayChannels[channelID] = types.Empty{}
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully marked channel as play channel!")
}

func (b *EoD) setModRole(roleID string, msg types.Msg, rsp types.Rsp) {
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

	lock.RLock()
	dat, exists := b.dat[msg.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	dat.ModRole = roleID
	lock.Lock()
	b.dat[msg.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully updated mod role!")
}

func (b *EoD) setUserColor(color string, removeColor bool, m types.Msg, rsp types.Rsp) {
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", m.GuildID, types.UserColor, m.Author.ID)
	var cnt int
	err := row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}

	// Remove color
	if cnt == 1 && removeColor {
		_, err = b.db.Exec("DELETE FROM eod_serverdata WHERE guild=? AND type=? AND value1=?", m.GuildID, types.PlayChannel, m.Author.ID)
		if rsp.Error(err) {
			return
		}

		lock.RLock()
		dat, exists := b.dat[m.GuildID]
		lock.RUnlock()
		if !exists {
			rsp.ErrorMessage("Guild not set up yet!")
			return
		}
		delete(dat.UserColors, m.Author.ID)
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()

		rsp.Resp("Successfully reset color.")
		return
	}

	if removeColor {
		rsp.ErrorMessage("You don't have a color!!")
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
		_, err = b.db.Exec("UPDATE eod_serverdata SET intval=? WHERE guild=? AND type=? AND value1=?", int(col), m.GuildID, types.ModRole, m.Author.ID)
		if rsp.Error(err) {
			return
		}
	} else {
		_, err = b.db.Exec("INSERT INTO eod_serverdata VALUES ( ?, ?, ?, ? )", m.GuildID, types.ModRole, m.Author.ID, int(col))
		if rsp.Error(err) {
			return
		}
	}

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		dat = types.NewServerData()
	}
	if dat.PlayChannels == nil {
		dat.PlayChannels = make(map[string]types.Empty)
	}
	dat.UserColors[m.Author.ID] = int(col)
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()

	rsp.Resp("Succesfully marked channel as play channel!")
}
