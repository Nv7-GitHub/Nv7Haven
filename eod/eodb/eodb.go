package eodb

import (
	"os"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type DB struct {
	sync.RWMutex

	Guild  string
	dbPath string

	Elements  []types.Element
	elemNames map[string]int
	combos    map[string]int                  // map["1+1"] = 5 for air + air = wind
	invs      map[string]*types.ElemContainer // map[userid]map[elemid]
	cats      map[string]*types.Category      // map[name]cat(id: cat name)
	Polls     map[string]types.Poll           // map[messageid]poll
	config    *types.ServerConfig

	invFiles   map[string]*os.File
	catFiles   map[string]*os.File
	elemFile   *os.File
	comboFile  *os.File
	configFile *os.File
}

func newDB(path string, guild string) *DB {
	return &DB{
		Guild:  guild,
		dbPath: path,

		combos:    make(map[string]int),
		invs:      make(map[string]*types.ElemContainer),
		cats:      make(map[string]*types.Category),
		Polls:     make(map[string]types.Poll),
		Elements:  make([]types.Element, 0),
		elemNames: make(map[string]int),

		invFiles: make(map[string]*os.File),
		catFiles: make(map[string]*os.File),
	}
}

func NewDB(guild, path string) (*DB, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	db := newDB(path, guild)

	// load
	err = db.loadElements()
	if err != nil {
		return nil, err
	}
	err = db.loadCombos()
	if err != nil {
		return nil, err
	}
	err = db.loadConfig()
	if err != nil {
		return nil, err
	}
	err = db.loadInvs()
	if err != nil {
		return nil, err
	}
	err = db.loadCats()
	if err != nil {
		return nil, err
	}
	err = db.loadPolls()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Close() {
	for _, file := range d.invFiles {
		file.Close()
	}
	for _, file := range d.catFiles {
		file.Close()
	}
	d.elemFile.Close()
	d.comboFile.Close()
	d.configFile.Close()
}
