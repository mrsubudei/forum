package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

func New() (*Sqlite, error) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		DB: db,
	}, nil
}

func (s *Sqlite) Close() {
	s.DB.Close()
}
