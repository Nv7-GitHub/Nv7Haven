package base

import (
	"log"

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
		CombCache:        make(map[string]types.CombCache),
		CommandStatsTODO: make(map[string]int),
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

const commandStatUpdateTrigger = 1000

func (b *Base) SaveCommandStats(guild string, mem *types.ServerMem) {
	if mem == nil {
		for gld, mem := range b.mem {
			b.SaveCommandStats(gld, mem)
		}
		return
	}

	todo := mem.CommandStatsTODO

	mem.Lock()
	mem.CommandStatsTODOCnt = 0
	mem.CommandStatsTODO = make(map[string]int)
	mem.Unlock()

	mem.RLock()
	for k, v := range todo {
		if v == 0 {
			continue
		}
		_, err := b.db.Exec("INSERT INTO command_stats (guild, command, count) VALUES ($1, $2, $3) ON CONFLICT (guild, command) DO UPDATE SET count = command_stats.count + $3", guild, k, v)
		if err != nil {
			log.Println("command stats write error", err)
		}
	}
	mem.RUnlock()
}

func (b *Base) IncrementCommandStat(c sevcord.Ctx, name string) {
	// Update command stats TODO
	mem := b.getMem(c)
	mem.Lock()
	mem.CommandStatsTODO[name]++
	mem.CommandStatsTODOCnt++
	mem.Unlock()

	if mem.CommandStatsTODOCnt >= commandStatUpdateTrigger {
		b.SaveCommandStats(c.Guild(), mem)
	}
}
