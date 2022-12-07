package sqlite_test

import (
	"reflect"
	"strings"
	"testing"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

func TestUserSqlite_Store(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal(err)
		}
		repo := sqlite.NewUsersRepo(db)

		user := entity.User{
			Id:    1,
			Name:  "Riddle",
			Email: "Riddle@mail.ru",
			Role:  "Пользователь",
			Sign:  " ",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal(err)
		}

		if id, err := repo.GetId(user); err != nil {
			t.Fatal(err)
		} else if id != int64(1) {
			t.Fatalf("ID=%v, want %v", id, int64(1))
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(user, found) {
			t.Fatalf("mismatch: %#v != %#v", user, found)
		}

		user2 := entity.User{
			Name: "buch",
		}
		if err := repo.Store(user2); err != nil {
			t.Fatal(err)
		}

		if id2, err := repo.GetId(user2); err != nil {
			t.Fatal(err)
		} else if id2 != int64(2) {
			t.Fatalf("ID=%v, want %v", id2, int64(2))
		}

	})

	t.Run("ErrNameAlreadyExist", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal(err)
		}
		repo := sqlite.NewUsersRepo(db)

		user := entity.User{
			Name:  "Riddle",
			Email: "Riddle@mail.ru",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal(err)
		}

		user2 := entity.User{
			Name: "Riddle",
		}
		if err = repo.Store(user2); err == nil {
			t.Fatalf("Error expected")
		} else if !strings.Contains(err.Error(), "UNIQUE constraint failed: users.name") {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}
