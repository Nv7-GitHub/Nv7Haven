package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) handlePollSuccess(p types.Poll) {
	b.lock.RLock()
	dat, exists := b.dat[p.Guild]
	b.lock.RUnlock()
	if !exists {
		return
	}

	controversial := dat.VoteCount != 0 && float32(p.Downvotes)/float32(dat.VoteCount) >= 0.3
	controversialTxt := ""
	if controversial {
		controversialTxt = " 🌩️"
	}

	switch p.Kind {
	case types.PollCombo:
		els, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			els = make([]string, len(dat))
			for i, val := range dat {
				els[i] = val.(string)
			}
		}
		b.elemCreate(p.Value3, els, p.Value4, controversialTxt, p.Guild)
	case types.PollSign:
		b.mark(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	case types.PollImage:
		b.image(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	case types.PollCategorize:
		els, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			els := make([]string, len(dat))
			for i, val := range dat {
				els[i] = val.(string)
			}
		}
		for _, val := range els {
			b.Categorize(val, p.Value1, p.Guild)
		}
		if len(els) == 1 {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Added **%s** to **%s** (By <@%s>)%s", els[0], p.Value1, p.Value4, controversialTxt))
		} else {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Added **%d elements** to **%s** (By <@%s>)%s", len(els), p.Value1, p.Value4, controversialTxt))
		}
	case types.PollUnCategorize:
		els := p.Data["elems"].([]string)
		for _, val := range els {
			b.UnCategorize(val, p.Value1, p.Guild)
		}
		if len(els) == 1 {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Removed **%s** from **%s** (By <@%s>)%s", els[0], p.Value1, p.Value4, controversialTxt))
		} else {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Removed **%d elements** from **%s** (By <@%s>)%s", len(els), p.Value1, p.Value4, controversialTxt))
		}
	case types.PollCatImage:
		b.catImage(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	case types.PollColor:
		b.color(p.Guild, p.Value1, p.Data["color"].(int), p.Value4, controversialTxt)
	case types.PollCatColor:
		b.catColor(p.Guild, p.Value1, p.Data["color"].(int), p.Value4, controversialTxt)
	}
}
