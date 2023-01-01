package base

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Base) getMem(c sevcord.Ctx) *types.ServerMem {
	b.lock.RLock()
	v, exists := b.mem[c.Guild()]
	b.lock.RUnlock()
	if exists {
		return v
	}

	v = &types.ServerMem{
		CombCache: make(map[string]types.CombCache),
	}

	b.lock.Lock()
	b.mem[c.Guild()] = v
	b.lock.Unlock()

	return v
}

func (b *Base) SaveCombCache(c sevcord.Ctx, comb types.CombCache) {
	mem := b.getMem(c)
	mem.Lock()
	mem.CombCache[c.Author().User.ID] = comb
	mem.Unlock()
}

func (b *Base) GetCombCache(c sevcord.Ctx) (types.CombCache, types.Resp) {
	mem := b.getMem(c)
	mem.RLock()
	v, exists := mem.CombCache[c.Author().User.ID]
	mem.RUnlock()
	if exists {
		return v, types.Ok()
	}
	return types.CombCache{}, types.Fail("You haven't combined anything!")
}
