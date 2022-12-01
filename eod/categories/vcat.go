package categories

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
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

	// Check if there already is an all elements vcat
	db.RLock()
	for _, vcat := range db.VCats() {
		if vcat.Rule == types.VirtualCategoryRuleAllElements {
			rsp.ErrorMessage(fmt.Sprintf("VCat **%s** already exists with the same rule!", vcat.Name)) // TODO: Translate
			db.RUnlock()
			return
		}
	}
	db.RUnlock()

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	if len(url.PathEscape(name)) > 1024 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
		return
	}
	vcat = &types.VirtualCategory{
		Name:      name,
		Guild:     m.GuildID,
		Creator:   m.Author.ID,
		Rule:      types.VirtualCategoryRuleAllElements,
		Data:      make(types.VirtualCategoryData),
		CreatedOn: types.NewTimeStamp(time.Now()),
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
	if len(url.PathEscape(name)) > 1024 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
		return
	}
	vcat = &types.VirtualCategory{
		Name:      name,
		Guild:     m.GuildID,
		Rule:      types.VirtualCategoryRuleRegex,
		Creator:   m.Author.ID,
		CreatedOn: types.NewTimeStamp(time.Now()),
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

	// Check if there already is an inv filter vcat
	db.RLock()
	for _, vcat := range db.VCats() {
		if vcat.Rule == types.VirtualCategoryRuleInvFilter && vcat.Data["user"] == user && vcat.Data["filter"] == filter {
			rsp.ErrorMessage(fmt.Sprintf("VCat **%s** already exists with the same rule!", vcat.Name)) // TODO: Translate
			db.RUnlock()
			return
		}
	}
	db.RUnlock()

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	if len(url.PathEscape(name)) > 1024 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
		return
	}
	vcat = &types.VirtualCategory{
		Name:      name,
		Guild:     m.GuildID,
		Creator:   m.Author.ID,
		CreatedOn: types.NewTimeStamp(time.Now()),
		Rule:      types.VirtualCategoryRuleInvFilter,
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
		rsp.ErrorMessage(db.Config.LangProperty("CategoryUsedInVCat", map[string]any{
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
		CreatedOn: types.NewTimeStamp(time.Now()),
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

	// Check if vcat already exists
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

	// Check if exists
	cat, res = db.GetCat(lhs)
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

	// Check if there already is an equivalent vcat
	db.RLock()
	for _, vcat := range db.VCats() {
		if vcat.Rule == types.VirtualCategoryRuleSetOperation && vcat.Data["lhs"] == lhs && vcat.Data["rhs"] == rhs && vcat.Data["operation"] == string(op) {
			rsp.ErrorMessage(fmt.Sprintf("VCat **%s** already exists with the same rule!", vcat.Name)) // TODO: Translate
			db.RUnlock()
			return
		}
	}
	db.RUnlock()

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	if len(url.PathEscape(name)) > 1024 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
		return
	}
	vcat = &types.VirtualCategory{
		Name:      name,
		Guild:     m.GuildID,
		Creator:   m.Author.ID,
		Rule:      types.VirtualCategoryRuleSetOperation,
		CreatedOn: types.NewTimeStamp(time.Now()),
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
	base.VCatCache = true
	for _, db := range b.DB {
		for _, cat := range db.VCats() {
			_, res := b.base.CalcVCat(cat, db, true)
			if !res.Exists { // Delete if not exists
				err := db.DeleteVCat(cat.Name)
				if err != nil {
					fmt.Println("VCat Calc Error: ", err)
				}
			}
		}
	}
}

func (b *Categories) VCatCreateInvhint(name string, elemName string, m types.Msg, rsp types.Rsp) {
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

	// Get elem
	elem, res := db.GetElementByName(elemName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Check if there already is an invhint vcat
	db.RLock()
	for _, vcat := range db.VCats() {
		if vcat.Rule == types.VirtualCategoryRuleInvhint && vcat.Data["element"] == float64(elem.ID) {
			rsp.ErrorMessage(fmt.Sprintf("VCat **%s** already exists with the same rule!", vcat.Name)) // TODO: Translate
			db.RUnlock()
			return
		}
	}
	db.RUnlock()

	// Create
	if strings.ToLower(name) == name {
		name = util.ToTitle(name)
	}
	if len(url.PathEscape(name)) > 1024 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
		return
	}
	vcat = &types.VirtualCategory{
		Name:      name,
		Guild:     m.GuildID,
		Rule:      types.VirtualCategoryRuleInvhint,
		Creator:   m.Author.ID,
		CreatedOn: types.NewTimeStamp(time.Now()),
		Data: types.VirtualCategoryData{
			"element": float64(elem.ID),
		},
	}
	err := db.SaveVCat(vcat)
	if rsp.Error(err) {
		return
	}
	rsp.Message("Created Virtual Category!") // TODO: Translate
}
