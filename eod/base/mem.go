package base

import (
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Base) getMem(c sevcord.Ctx) *types.ServerMem {
	return b.getMemGuild(c.Guild())
}
func (b *Base) getMemGuild(guild string) *types.ServerMem {
	b.lock.RLock()
	v, exists := b.mem[guild]
	b.lock.RUnlock()
	if exists {
		return v
	}

	v = &types.ServerMem{
		CombCache:        make(map[string]types.CombCache),
		CommandStatsTODO: make(map[string]int),
	}

	b.lock.Lock()
	b.mem[guild] = v
	b.lock.Unlock()

	return v
}
func (b *Base) SaveCombCacheUser(user string, guild string, comb types.CombCache) {
	mem := b.getMemGuild(guild)
	mem.Lock()
	mem.CombCache[user] = comb
	mem.Unlock()
}
func (b *Base) SaveCombCache(c sevcord.Ctx, comb types.CombCache) {

	b.SaveCombCacheUser(c.Author().User.ID, c.Guild(), comb)
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
