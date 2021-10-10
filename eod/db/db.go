package db

import "database/sql"

type DB struct {
	db    *sql.DB
	cache map[string]*sql.Stmt
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	quer, err := d.getStmt(query)
	if err != nil {
		return nil, err
	}
	return quer.Exec(args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	quer, err := d.getStmt(query)
	if err != nil {
		return nil, err
	}
	return quer.Query(args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	quer, err := d.getStmt(query)
	if err != nil {
		return d.db.QueryRow(query, args...) // Can't return an error, so have the DB do it instead
	}
	return quer.QueryRow(args...)
}

func (d *DB) getStmt(query string) (*sql.Stmt, error) {
	quer, exists := d.cache[query]
	if !exists {
		stmt, err := d.db.Prepare(query)
		if err != nil {
			return nil, err
		}
		d.cache[query] = stmt
		quer = stmt
	}
	return quer, nil
}

func NewDB(db *sql.DB) *DB {
	return &DB{db: db, cache: make(map[string]*sql.Stmt)}
}

func (d *DB) GetSqlDB() *sql.DB {
	return d.db
}

func (d *DB) Close() {
	for _, stmt := range d.cache {
		stmt.Close()
	}
}
