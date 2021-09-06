package anarchy

import (
	"database/sql"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
)

var lock *sync.RWMutex

type Anarchy struct {
	*pb.UnimplementedAnarchyServer

	db      *sql.DB
	cache   map[string]*pb.AnarchyElement
	recents *sync.Cond
}

func (a *Anarchy) init() {
	var cnt int
	err := a.db.QueryRow(`SELECT COUNT(1) FROM anarchy_elements`).Scan(&cnt)
	if err != nil {
		panic(err)
	}

	bar := progressbar.New(cnt)

	res, err := a.db.Query("SELECT * FROM anarchy_elements WHERE 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	for res.Next() {
		elem := &pb.AnarchyElement{}
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.CreatedOn, &elem.Complexity, &elem.Uses, &elem.FoundBy)
		if err != nil {
			panic(err)
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		a.cache[elem.Name] = elem

		bar.Add(1)
	}
	bar.Finish()
}

// InitAnarchy initializes Elemental 7's Anarchy server
func InitAnarchy(db *sql.DB, grpc *grpc.Server) *Anarchy {
	a := &Anarchy{
		db:      db,
		cache:   make(map[string]*pb.AnarchyElement),
		recents: sync.NewCond(&sync.Mutex{}),
	}
	a.init()

	pb.RegisterAnarchyServer(grpc, a)

	return a
}
