package sqlite_test

// import (
// 	"reflect"
// 	"testing"

// 	"forum/internal/entity"
// 	"forum/internal/repository/sqlite"
// 	"forum/pkg/sqlite3"
// )

// func TestUserSqlite_Store(t *testing.T) {
// 	t.Run("OK", func(t *testing.T) {
// 		db := sqlite3.Sqlite
// 		defer MustCloseDB(t, db)

// 		repo := sqlite.NewUsersRepo(db)

// 		user := entity.User{
// 			Name:  "Riddle",
// 			Email: "Riddle@mail.ru",
// 		}

// 		if err := repo.Store(user); err != nil {
// 			t.Fatal(err)
// 		} else if got, want := user.Id, int64(1); got != want {
// 			t.Fatalf("ID=%v, want %v", got, want)
// 		}

// 		if found, err := repo.GetById(1); err != nil {
// 			t.Fatal(err)
// 		} else if !reflect.DeepEqual(user, found) {
// 			t.Fatalf("mismatch: %#v != %#v", user, found)
// 		}
// 	})
// }
