package polls

import (
	"database/sql"
	"log"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const UpArrow = "⬆️"
const DownArrow = "⬇️"
const TriageDelay = time.Second

func (b *Polls) reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	if r.Emoji.Name != UpArrow && r.Emoji.Name != DownArrow {
		return
	}

	// Get poll & vote cnt
	var p types.Poll
	err := b.db.Get(&p, "SELECT * FROM polls WHERE guild=$1 AND message=$2", r.GuildID, r.MessageID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("poll fetch err", err)
		}
		return
	}
	var votecnt int
	err = b.db.QueryRow("SELECT votecnt FROM config WHERE guild=$1", r.GuildID).Scan(&votecnt)
	if err != nil {
		log.Println("poll votecnt err", err)
		return
	}

	// Triage polls, once per guild
	if b.shouldTriage(r.GuildID) {
		go b.triagePolls(r.GuildID, s)
	}

	// User trying to delete?
	if r.UserID == p.Creator && r.Emoji.Name == DownArrow {
		b.deletePoll(&p, s)
		return
	}

	// Update user vote count
	_, err = b.db.Exec(`UPDATE inventories SET votecnt=votecnt+1 WHERE guild=$1 AND "user"=$2`, p.Guild, r.UserID)
	if err != nil {
		log.Println("user votecnt update err", err)
	}

	// Handle
	if r.Emoji.Name == UpArrow {
		p.Upvotes++
	} else {
		p.Downvotes++
	}

	// Update
	_, err = b.db.NamedExec("UPDATE polls SET upvotes=:upvotes, downvotes=:downvotes WHERE guild=:guild AND message=:message", p)
	if err != nil {
		log.Println("poll update err", err)
		return
	}

	// Check
	b.checkPoll(&p, votecnt, s)
}

func (b *Polls) unReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}
	if r.Emoji.Name != UpArrow && r.Emoji.Name != DownArrow {
		return
	}

	// Get poll & vote cnt
	var p types.Poll
	err := b.db.Get(&p, "SELECT * FROM polls WHERE guild=$1 AND message=$2", r.GuildID, r.MessageID)
	if err != nil {
		return
	}
	var votecnt int
	err = b.db.QueryRow("SELECT votecnt FROM config WHERE guild=$1", r.GuildID).Scan(&votecnt)
	if err != nil {
		return
	}

	// Triage polls, once per guild
	if b.shouldTriage(r.GuildID) {
		go b.triagePolls(r.GuildID, s)
	}

	// Update user vote count
	_, err = b.db.Exec(`UPDATE inventories SET votecnt=votecnt-1 WHERE guild=$1 AND "user"=$2`, p.Guild, r.UserID)
	if err != nil {
		log.Println("user votecnt update err", err)
	}

	// Handle
	if r.Emoji.Name == UpArrow {
		p.Upvotes--
	} else {
		p.Downvotes--
	}

	// Update
	_, err = b.db.NamedExec("UPDATE polls SET upvotes=:upvotes, downvotes=:downvotes WHERE guild=:guild AND message=:message", p)
	if err != nil {
		return
	}

	// Check
	b.checkPoll(&p, votecnt, s)
}

func (b *Polls) shouldTriage(guild string) bool {
	b.lock.RLock()
	triaged := b.triaged[guild]
	b.lock.RUnlock()
	if triaged {
		return false
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	if b.triaged[guild] { // Another event may have triaged it
		return false
	}
	b.triaged[guild] = true
	return true
}

func (b *Polls) triagePolls(guild string, dg *discordgo.Session) {
	// Get polls & vote cnt
	var votecnt int
	err := b.db.QueryRow("SELECT votecnt FROM config WHERE guild=$1", guild).Scan(&votecnt)
	if err != nil {
		log.Println("triage votecnt err", err)
		return
	}
	var polls []types.Poll
	err = b.db.Select(&polls, "SELECT * FROM polls WHERE guild=$1", guild)
	if err != nil {
		log.Println("triage fetch err", err)
		return
	}

	// Handle
	for i := range polls {
		time.Sleep(TriageDelay) // Go slowly, this isn't urgent
		b.triagePoll(&polls[i], votecnt, dg)
	}
}

func (b *Polls) triagePoll(p *types.Poll, votecnt int, dg *discordgo.Session) {
	// Old polls may have data this version can't handle
	defer func() {
		if r := recover(); r != nil {
			log.Println("triage panic", p.Guild, p.Message, r)
		}
	}()

	// Poll may have been voted on while triaging
	err := b.db.Get(p, "SELECT * FROM polls WHERE guild=$1 AND message=$2", p.Guild, p.Message)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("triage poll err", err)
		}
		return
	}

	// Get message
	msg, err := dg.ChannelMessage(p.Channel, p.Message)
	if err != nil {
		// Message deleted?
		rest, ok := err.(*discordgo.RESTError)
		if ok && rest.Message != nil && rest.Message.Code == discordgo.ErrCodeUnknownMessage {
			b.pollReject(p, dg)
			return
		}
		log.Println("triage message err", err)
		return
	}

	// Recount votes
	p.Upvotes = 0
	p.Downvotes = 0
	for _, react := range msg.Reactions {
		cnt := react.Count
		if react.Me { // The bot's own reaction doesn't count
			cnt--
		}
		switch react.Emoji.Name {
		case UpArrow:
			p.Upvotes = cnt

		case DownArrow:
			p.Downvotes = cnt
		}
	}

	// User trying to delete?
	if p.Downvotes > 0 {
		users, err := dg.MessageReactions(p.Channel, p.Message, DownArrow, 100, "", "")
		if err != nil {
			log.Println("triage reactions err", err)
			return
		}
		for _, user := range users {
			if user.ID == p.Creator {
				b.deletePoll(p, dg)
				return
			}
		}
	}

	// Update
	_, err = b.db.Exec("UPDATE polls SET upvotes=$1, downvotes=$2 WHERE guild=$3 AND message=$4", p.Upvotes, p.Downvotes, p.Guild, p.Message)
	if err != nil {
		log.Println("triage update err", err)
		return
	}

	// Check
	b.checkPoll(p, votecnt, dg)
}
