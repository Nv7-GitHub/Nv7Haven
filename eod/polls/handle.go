package polls

import (
	"database/sql"
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const UpArrow = "⬆️"
const DownArrow = "⬇️"

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

	// User trying to delete?
	if r.UserID == p.Creator && r.Emoji.Name == DownArrow {
		b.deletePoll(&p, s)
		return
	}

	// Update vote count
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
