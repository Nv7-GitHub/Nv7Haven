package categories

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Categories) VCatCreateAllElementsCmd(name string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	// Check if exists
	cat, res := db.GetCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", cat.Name))
		return
	}
	vcat, res := db.GetVCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", vcat.Name))
		return
	}

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	vcat = &types.VirtualCategory{
		Name:  name,
		Guild: m.GuildID,
		Rule:  types.VirtualCategoryRuleAllElements,
		Data:  make(types.VirtualCategoryData),
	}
	err := db.SaveVCat(vcat)
	if rsp.Error(err) {
		return
	}
	rsp.Message("Created Virtual Category!") // TODO: Translate
}

func (b *Categories) VCatCreateRegexCmd(name string, regex string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	// Check if exists
	cat, res := db.GetCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", cat.Name))
		return
	}
	vcat, res := db.GetVCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", vcat.Name))
		return
	}

	// Check if valid regex
	_, err := regexp.Compile(regex)
	if rsp.Error(err) {
		return
	}

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	vcat = &types.VirtualCategory{
		Name:  name,
		Guild: m.GuildID,
		Rule:  types.VirtualCategoryRuleRegex,
		Data: types.VirtualCategoryData{
			"regex": regex,
		},
	}
	err = db.SaveVCat(vcat)
	if rsp.Error(err) {
		return
	}
	rsp.Message("Created Virtual Category!") // TODO: Translate
}

func (b *Categories) VCatCreateInvFilterCmd(name string, user string, filter string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	// Check if exists
	cat, res := db.GetCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", cat.Name))
		return
	}
	vcat, res := db.GetVCat(name)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", vcat.Name))
		return
	}

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	vcat = &types.VirtualCategory{
		Name:  name,
		Guild: m.GuildID,
		Rule:  types.VirtualCategoryRuleInvFilter,
		Data: types.VirtualCategoryData{
			"user":   user,
			"filter": filter,
		},
	}
	err := db.SaveVCat(vcat)
	if rsp.Error(err) {
		return
	}
	rsp.Message("Created Virtual Category!") // TODO: Translate
}

func (b *Categories) DeleteVCatCmd(name string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	// Check if exists
	vcat, res := db.GetVCat(name)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Check if affects anything else
	ok := true
	var vname string
	db.RLock()
	for _, cat := range db.VCats() {
		if cat.Rule == types.VirtualCategoryRuleSetOperation {
			if strings.EqualFold(cat.Data["lhs"].(string), vcat.Name) {
				vname = cat.Name
				ok = false
				break
			}

			if strings.EqualFold(cat.Data["rhs"].(string), vcat.Name) {
				vname = cat.Name
				ok = false
				break
			}
		}
	}
	db.RUnlock()
	if !ok {
		rsp.ErrorMessage(db.Config.LangProperty("CategoryUsedInVCat", map[string]interface{}{
			"Category":        vcat.Name,
			"VirtualCategory": vname,
		}))
		return
	}

	err := b.polls.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollDeleteVCat,
		Suggestor: m.Author.ID,
		PollVCatDeleteData: &types.PollVCatDeleteData{
			Category: vcat.Name,
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(db.Config.LangProperty("VCatDeleteSuggested", vcat.Name))
}

func (b *Categories) VCatOpCmd(op types.CategoryOperation, name string, lhs string, rhs string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	// Check if exists
	cat, res := db.GetCat(lhs)
	if !res.Exists {
		vcat, res := db.GetVCat(lhs)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		} else {
			lhs = vcat.Name
		}
	} else {
		lhs = cat.Name
	}
	cat, res = db.GetCat(rhs)
	if !res.Exists {
		vcat, res := db.GetVCat(rhs)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		} else {
			rhs = vcat.Name
		}
	} else {
		rhs = cat.Name
	}

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	vcat := &types.VirtualCategory{
		Name:  name,
		Guild: m.GuildID,
		Rule:  types.VirtualCategoryRuleSetOperation,
		Data: types.VirtualCategoryData{
			"lhs":       lhs,
			"rhs":       rhs,
			"operation": string(op),
		},
	}
	err := db.SaveVCat(vcat)
	if rsp.Error(err) {
		return
	}
	rsp.Message("Created Virtual Category!") // TODO: Translate
}

func (b *Categories) CacheVCats() {
	for _, db := range b.DB {
		regs := make(map[string]*regexp.Regexp)
		vcats := make(map[string]map[int]types.Empty)

		// Performance
		times := make(map[string]time.Duration)
		code := make(map[string]string)
		for _, vcat := range db.VCats() {
			if vcat.Rule == types.VirtualCategoryRuleRegex {
				vcat.Cache = make(map[int]types.Empty)
				vcats[vcat.Name] = vcat.Cache
				regs[vcat.Name] = regexp.MustCompile(vcat.Data["regex"].(string))

				// Performance
				times[vcat.Name] = 0
				code[vcat.Name] = vcat.Data["regex"].(string)
			}
		}

		for _, el := range db.Elements {
			for k, reg := range regs {
				start := time.Now()
				match := reg.Match([]byte(el.Name))
				times[k] += time.Since(start)
				if match {
					vcats[k][el.ID] = types.Empty{}
				}
			}
		}

		// Performance
		type timeData struct {
			dur  time.Duration
			code string
		}
		i := 0
		dat := make([]timeData, len(times))
		for k, v := range times {
			dat[i] = timeData{v, code[k]}
			i++
		}
		sort.Slice(dat, func(i, j int) bool {
			return dat[i].dur > dat[j].dur
		})
		for _, v := range dat {
			fmt.Printf("%s: %s\n", v.dur.String(), v.code)
		}
	}
}
