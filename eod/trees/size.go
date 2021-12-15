package trees

import (
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type SizeTree struct {
	Size  int
	db    *eodb.DB
	added map[int]types.Empty
}

func (s *SizeTree) AddElem(elem int, notoplevel ...bool) (bool, string) {
	_, exists := s.added[elem]
	if exists {
		return true, ""
	}

	if len(notoplevel) == 0 {
		s.db.RLock()
		defer s.db.RUnlock()
	}

	el, res := s.db.GetElement(elem, true)
	if !res.Exists {
		return false, res.Message
	}

	for _, par := range el.Parents {
		suc, msg := s.AddElem(par, true)
		if !suc {
			return false, msg
		}
	}

	s.added[elem] = types.Empty{}
	s.Size++

	return true, ""
}

func NewSizeTree(db *eodb.DB) *SizeTree {
	return &SizeTree{Size: 0, db: db, added: make(map[int]types.Empty)}
}

func ElemCreateSize(parents []int, db *eodb.DB) (int, bool, string) {
	size := NewSizeTree(db)
	for _, par := range parents {
		suc, msg := size.AddElem(par)
		if !suc {
			return 0, false, msg
		}
	}
	return size.Size + 1, true, ""
}
