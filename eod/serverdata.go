package eod

import (
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
