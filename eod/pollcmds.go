package eod

func (b *EoD) suggestCmd(suggestion string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	comb, exists := dat.combCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined anything!")
		return
	}

	b.createPoll(poll{
		Channel:   dat.votingChannel,
		Guild:     m.GuildID,
		Kind:      pollCombo,
		Value1:    comb.elem1,
		Value2:    comb.elem2,
		Value3:    suggestion,
		Value4:    m.Author.ID,
		Data:      make(map[string]interface{}),
		Upvotes:   0,
		Downvotes: 0,
	})
}
