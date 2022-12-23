package sqlite_test

import (
	"database/sql"
	"errors"
	"testing"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

func TestCommentStore(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  1,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		comment2 := entity.Comment{
			PostId:  1,
			User:    entity.User{Id: 2},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment2); err != nil {
			t.Fatal("Unable to store:", err)
		}
	})
}

func TestCommentFetch(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  1,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		comment2 := entity.Comment{
			PostId:  1,
			User:    entity.User{Id: 2},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if found, err := repo.Fetch(1); err != nil {
			t.Fatal("Unable to Fetch:", err)
		} else if len(found) != 2 {
			t.Fatalf("want len = %d, got len = %d:", 2, len(found))
		}
	})
}

func TestCommentGetyId(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  7,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		comment2 := entity.Comment{
			PostId:  3,
			User:    entity.User{Id: 2},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment2); err != nil {
			t.Fatal("Unable to Store:", err)
		}

		if found, err := repo.GetById(2); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Id != 2 {
			t.Fatalf("want len = %d, got len = %d:", 2, found.Id)
		}
	})
}

func TestCommentUpdate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  7,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		newContent := "New Content"
		comment.Id = 1
		comment.Content = newContent

		if err = repo.Update(comment); err != nil {
			t.Fatal("Unable to Update:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Content != newContent {
			t.Fatalf("want = %v, got = %v:", newContent, found.Content)
		}
	})
}

func TestGetPostIds(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  7,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		user := entity.User{Id: 1}
		comment2 := entity.Comment{
			PostId:  2,
			User:    user,
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if found, err := repo.GetPostIds(user); err != nil {
			t.Fatal("Unable to GetPostIds:", err)
		} else if len(found) != 2 {
			t.Fatalf("want = %d, got = %d:", 2, len(found))
		}
	})
}

func TestCommentDelete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			PostId:  7,
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(comment); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if _, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		}
		comment.Id = 1

		if err := repo.Delete(comment); err != nil {
			t.Fatal("Unable to Delete:", err)
		}

		if _, err := repo.GetById(1); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("want err = %v, got err = %v:", sql.ErrNoRows, err)
		}
	})
}

func TestCommentStoreLike(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			Id:   1,
			User: entity.User{Id: 1},
			Date: "2022-19-01",
		}

		if err := repo.StoreLike(comment); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}
	})
}

func TestCommentDeleteLike(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			Id:   1,
			User: entity.User{Id: 1},
			Date: "2022-19-01",
		}

		if err = repo.StoreLike(comment); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		if err = repo.DeleteLike(comment); err != nil {
			t.Fatal("Unable to DeleteLike:", err)
		}
	})
}

func TestCommentStoreDislike(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			Id:   1,
			User: entity.User{Id: 1},
			Date: "2022-19-01",
		}

		if err := repo.StoreDislike(comment); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}
	})
}

func TestCommentDeleteDislike(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			Id:   1,
			User: entity.User{Id: 1},
			Date: "2022-19-01",
		}

		if err = repo.StoreDislike(comment); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		if err = repo.DeleteDislike(comment); err != nil {
			t.Fatal("Unable to DeleteDislike:", err)
		}
	})
}

func TestCommentFetchReactions(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewCommentsRepo(db)

		comment := entity.Comment{
			Id:   1,
			User: entity.User{Id: 1},
			Date: "2022-19-01",
		}

		if err = repo.StoreDislike(comment); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		comment.User.Id = 2
		if err = repo.StoreDislike(comment); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		comment.User.Id = 3
		if err = repo.StoreLike(comment); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		if found, err := repo.FetchReactions(1); err != nil {
			t.Fatal("Unable to FetchReactions:", err)
		} else if len(found.Likes) != 1 {
			t.Fatalf("want  = %d, got = %d:", 1, len(found.Likes))
		} else if len(found.Dislikes) != 2 {
			t.Fatalf("want  = %d, got = %d:", 2, len(found.Dislikes))
		}
	})
}
