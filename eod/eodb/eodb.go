package eodb

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type Data struct {
	*sync.RWMutex

	DB   map[string]*DB
	Data map[string]*types.ServerData
	path string
}

func NewData(path string) (*Data, error) {
	folders, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	d := &Data{
		RWMutex: &sync.RWMutex{},

		DB:   make(map[string]*DB),
		Data: make(map[string]*types.ServerData),
		path: path,
	}
	for _, folder := range folders {
		db, err := NewDB(folder.Name(), filepath.Join(path, folder.Name()))
		if err != nil {
			return nil, err
		}
		d.DB[folder.Name()] = db
		d.Data[folder.Name()] = types.NewServerData()
	}

	return d, nil
}

func (d *Data) GetDB(guild string) (*DB, types.GetResponse) {
	d.RLock()
	db, exists := d.DB[guild]
	d.RUnlock()
	if !exists {
		return nil, types.GetResponse{
			Exists:  false,
			Message: "Guild not found",
		}
	}
	return db, types.GetResponse{Exists: true}
}

func (d *Data) GetData(guild string) (*types.ServerData, types.GetResponse) {
	d.RLock()
	data, exists := d.Data[guild]
	d.RUnlock()
	if !exists {
		return nil, types.GetResponse{
			Exists:  false,
			Message: "Guild not found",
		}
	}
	return data, types.GetResponse{Exists: true}
}

func (d *Data) NewDB(guild string) (*DB, error) {
	d.Lock()
	defer d.Unlock()
	return NewDB(guild, filepath.Join(d.path, guild))
}

type DB struct {
	sync.RWMutex

	Guild  string
	dbPath string

	Elements  []types.Element
	elemNames map[string]int
	combos    map[string]int              // map["1+1"] = 5 for air + air = wind
	invs      map[string]*types.Inventory // map[userid]map[elemid]
	cats      map[string]*types.Category  // map[name]cat(id: cat name)
	Polls     map[string]types.Poll       // map[messageid]poll
	Config    *types.ServerConfig

	invFiles   map[string]*os.File
	catFiles   map[string]*os.File
	elemFile   *os.File
	comboFile  *os.File
	configFile *os.File
}

func (d *DB) Invs() map[string]*types.Inventory {
	return d.invs
}

func (d *DB) Cats() map[string]*types.Category {
	return d.cats
}

func (d *DB) ComboCnt() int {
	return len(d.combos)
}

func newDB(path string, guild string) *DB {
	return &DB{
		Guild:  guild,
		dbPath: path,

		combos:    make(map[string]int),
		invs:      make(map[string]*types.Inventory),
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
