package sqlite_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

func TestStore(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewUsersRepo(db)

		user := entity.User{
			Id:    1,
			Name:  "Riddle",
			Email: "Riddle@mail.ru",
			City:  "Astana",
			Role:  "Пользователь",
			Sign:  " ",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if id, err := repo.GetId(user); err != nil {
			t.Fatal("Unable to GetId:", err)
		} else if id != int64(1) {
			t.Fatalf("ID=%v, want %v", id, int64(1))
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if !reflect.DeepEqual(user, found) {
			t.Fatalf("mismatch: %#v != %#v", user, found)
		}

		user2 := entity.User{
			Name: "buch",
		}
		if err := repo.Store(user2); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		if id2, err := repo.GetId(user2); err != nil {
			t.Fatal("Unable to GetId:", err)
		} else if id2 != int64(2) {
			t.Fatalf("ID=%v, want %v", id2, int64(2))
		}
	})

	t.Run("ErrNameAlreadyExist", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}
		repo := sqlite.NewUsersRepo(db)

		user := entity.User{
			Name:  "Riddle",
			Email: "Riddle@mail.ru",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
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

func TestFetch(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}
		repo := sqlite.NewUsersRepo(db)

		user := entity.User{
			Name:  "Riddle",
			Email: "Riddle@mail.ru",
			City:  "Astana",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		if err := repo.Store(entity.User{Name: "adf"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		users, err := repo.Fetch()
		if err != nil {
			t.Fatal("Unable to Fetch:", err)
		}
		if len(users) != 2 {
			t.Fatalf("want %d length, got %d length", 2, len(users))
		}
	})
}

func TestGetId(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		if err := repo.Store(entity.User{Name: "qwe", Email: "hthth"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}
		if err := repo.Store(entity.User{Name: "sdf"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		user := entity.User{
			Name:  "Subi",
			Email: "Subi@mail.ru",
			City:  "Astana",
		}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}
		time := time.Now()
		u := entity.User{Id: 3, SessionToken: "newToken!@#$%", SessionTTL: time}
		if err = repo.NewSession(u); err != nil {
			t.Fatal("Unable to NewSession:", err)
		}

		userToFindByToken := entity.User{SessionToken: "newToken!@#$%"}
		if id, err := repo.GetId(userToFindByToken); err != nil {
			t.Fatal("Unable to GetId:", err)
		} else if id != 3 {
			t.Fatalf("want id = %d, got id = %d", 3, id)
		}

		userToFindByName := entity.User{Name: "Subi"}
		if id, err := repo.GetId(userToFindByName); err != nil {
			t.Fatal("Unable to GetId:", err)
		} else if id != 3 {
			t.Fatalf("want id = %d, got id = %d", 3, id)
		}

		userToFindByEmail := entity.User{Email: "Subi@mail.ru"}
		if id, err := repo.GetId(userToFindByEmail); err != nil {
			t.Fatal("Unable to GetId:", err)
		} else if id != 3 {
			t.Fatalf("want id = %d, got id = %d", 3, id)
		}
	})
}

func TestGetById(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		if err := repo.Store(entity.User{Name: "Bobik", Email: "hthth"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}
		if err := repo.Store(entity.User{Name: "Tuzik"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		if foundUser, err := repo.GetById(2); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if foundUser.Name != "Tuzik" {
			t.Fatalf("want name = Tuzik, got name = %v", foundUser.Name)
		}
	})
}

func TestGetSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		if err := repo.Store(entity.User{Name: "Bobik", Email: "hthth"}); err != nil {
			t.Fatal("Unable to Store:", err)
		}
		session := "Token!@#$%^&"
		time := time.Now()
		if err = repo.NewSession(entity.User{
			Id: 1, SessionToken: session,
			SessionTTL: time,
		}); err != nil {
			t.Fatal("Unable to NewSession:", err)
		}

		if foundUser, err := repo.GetSession(1); err != nil {
			t.Fatal("Unable to GetSession:", err)
		} else if foundUser.SessionToken != "Token!@#$%^&" {
			t.Fatalf("want session = %v, got session = %v", session, foundUser.SessionToken)
		} else if foundUser.SessionTTL.Format("2006-01-02 15:04:05") != time.Format("2006-01-02 15:04:05") {
			t.Fatalf("want TTL = %v, got TTL = %v", time, foundUser.SessionTTL)
		}
	})
}

func TestUpdateInfo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		user := entity.User{Name: "Bobik", Email: "hthth@dfg"}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		user.Id = 1
		user.City = "Astana"
		user.Role = "Moderator"

		if err != repo.UpdateInfo(user) {
			t.Fatal("Unable to UpdateInfo:", err)
		}

		if foundUser, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if foundUser.City != user.City {
			t.Fatalf("want city = %v, got city = %v", user.City, foundUser.City)
		} else if foundUser.Role != user.Role {
			t.Fatalf("want role = %v, got role = %v", user.Role, foundUser.Role)
		}
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		user := entity.User{Name: "Bobik", Email: "hthth@dfg", Password: "password!@#$%^"}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		user.Id = 1
		user.Password = "NewPass!@#"

		if err != repo.UpdatePassword(user) {
			t.Fatal("Unable to UpdatePassword:", err)
		}

		if foundUser, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if foundUser.Password != user.Password {
			t.Fatalf("want password = %v, got password = %v", user.Password, foundUser.Password)
		}
	})
}

func TestUpdateSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		user := entity.User{Name: "Bobik", Email: "hthth@dfg", Password: "password!@#$%^"}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		timer := time.Now()
		u := entity.User{Id: 1, SessionToken: "newToken!@#$%", SessionTTL: timer}
		if err = repo.NewSession(u); err != nil {
			t.Fatal("Unable to NewSession:", err)
		}

		u.SessionTTL = timer.Add(time.Hour)

		if err = repo.UpdateSession(u); err != nil {
			t.Fatal("Unable to UpdateSession:", err)
		}

		if foundUser, err := repo.GetSession(1); err != nil {
			t.Fatal("Unable to GetSession:", err)
		} else if foundUser.SessionTTL.Format("2006-01-02 15:04:05") != u.SessionTTL.Format("2006-01-02 15:04:05") {
			t.Fatalf("want TTL = %v, got TTL = %v", u.SessionTTL, foundUser.SessionTTL)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to CreateDB:", err)
		}

		repo := sqlite.NewUsersRepo(db)

		user := entity.User{Name: "Bobik", Email: "hthth@dfg", Password: "password!@#$%^"}

		if err := repo.Store(user); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		if err = repo.Delete(entity.User{Id: 1}); err != nil {
			t.Fatal("Unable to Delete:", err)
		}

		if _, err := repo.GetById(1); err == nil {
			t.Fatal("Expected error:", err)
		}
	})
}
