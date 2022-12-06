package sqlite3_test

import (
	"testing"

	"forum/pkg/sqlite3"
)

func TestNew(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

func MustOpenDB(tb testing.TB) *sqlite3.Sqlite {
	tb.Helper()

	db, err := sqlite3.New()
	if err != nil {
		tb.Fatal(err)
	}
	return db
}

func MustCloseDB(tb testing.TB, db *sqlite3.Sqlite) {
	tb.Helper()
	if err := db.DB.Close(); err != nil {
		tb.Fatal(err)
	}
}
