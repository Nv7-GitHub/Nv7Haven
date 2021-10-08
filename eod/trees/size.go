package trees

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type SizeTree struct {
	Size  int
	dat   types.ServerData
	added map[string]types.Empty
}

func (s *SizeTree) AddElem(name string, notoplevel ...bool) (bool, string) {
	_, exists := s.added[name]
	if exists {
		return true, ""
	}

	if len(notoplevel) == 0 {
		s.dat.Lock.RLock()
		defer s.dat.Lock.RUnlock()
	}

	el, res := s.dat.GetElement(name, true)
	if !res.Exists {
		return false, res.Message
	}

	for _, par := range el.Parents {
		suc, msg := s.AddElem(par, true)
		if !suc {
			return false, msg
		}
	}

	s.added[strings.ToLower(name)] = types.Empty{}
	s.Size++

	return true, ""
}

func NewSizeTree(dat types.ServerData) *SizeTree {
	return &SizeTree{Size: 0, dat: dat, added: make(map[string]types.Empty)}
}

func ElemCreateSize(parents []string, dat types.ServerData) (int, bool, string) {
	size := NewSizeTree(dat)
	for _, par := range parents {
		suc, msg := size.AddElem(par)
		if !suc {
			return 0, false, msg
		}
	}
	return size.Size + 1, true, ""
}
