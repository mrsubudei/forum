package sqlite

import (
	"testing"

	"forum/pkg/sqlite3"
)

func TestNew(t *testing.T) {
	db := MustOpenDB(t, "../../../database/forum.db")
	MustCloseDB(t, db)
}

func MustOpenDB(tb testing.TB, path string) *sqlite3.Sqlite {
	db, err := sqlite3.New(path)
	if err != nil {
		tb.Fatal(err)
	}
	return db
}

func MustCloseDB(tb testing.TB, db *sqlite3.Sqlite) {
	if err := db.DB.Close(); err != nil {
		tb.Fatal(err)
	}
}
