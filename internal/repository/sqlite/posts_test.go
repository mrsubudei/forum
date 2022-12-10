package sqlite_test

import (
	"reflect"
	"strings"
	"testing"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

func TestPostStore(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		post := entity.Post{
			User:    entity.User{Id: 1},
			Date:    "2022-19-01",
			Title:   "Cars",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post); err != nil {
			t.Fatal("Unable to store:", err)
		} else if post.Id != 1 {
			t.Fatalf("want id = %d, got id = %d:", post.Id, 1)
		}

		post2 := entity.Post{
			User:    entity.User{Id: 2},
			Date:    "2022-19-02",
			Title:   "Sports",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		} else if post2.Id != 2 {
			t.Fatalf("want id = %d, got id = %d:", post2.Id, 2)
		}
	})
}

func TestStoreTopicReference(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		categories := []string{"cars", "cinema", "food"}

		post := entity.Post{Id: 1, Categories: categories}

		if err := repo.StoreTopicReference(post); err != nil {
			t.Fatal("Unable to StoreTopicReference:", err)
		}
	})
}

func TestPostFetch(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		post := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Title:   "Cars",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post); err != nil {
			t.Fatal("Unable to store:", err)
		} else if post.Id != 1 {
			t.Fatalf("want id = %d, got id = %d:", post.Id, 1)
		}

		post2 := entity.Post{
			User:    entity.User{Id: 2, Name: "Subi"},
			Date:    "2022-19-02",
			Title:   "Sports",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		} else if post2.Id != 2 {
			t.Fatalf("want id = %d, got id = %d:", post2.Id, 2)
		}

		if fetched, err := repo.Fetch(); err != nil {
			t.Fatal("Unable to Fetch:", err)
		} else if len(fetched) != 2 {
			t.Fatalf("want len = %d, got len = %d:", len(fetched), 2)
		}
	})
}

func TestFetchByAuthor(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		post := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Title:   "Cars",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post); err != nil {
			t.Fatal("Unable to store:", err)
		}

		title := "Travel"
		post1 := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Title:   title,
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		post2 := entity.Post{
			User:    entity.User{Id: 2, Name: "Subi"},
			Date:    "2022-19-02",
			Title:   "Sports",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if fetched, err := repo.FetchByAuthor(entity.User{Id: 1}); err != nil {
			t.Fatal("Unable to FetchByAuthor:", err)
		} else if len(fetched) != 2 {
			t.Fatalf("want len = %d, got len = %d:", len(fetched), 2)
		} else if fetched[1].Title != title {
			t.Fatalf("want title = %v, got title = %v:", title, fetched[1].Title)
		}
	})
}

func TestFetchIdsByReaction(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		user := entity.User{Id: 1, Name: "Riddle"}
		post := entity.Post{
			User: user,
			Date: "2022-19-01",
		}

		if err := repo.Store(&post); err != nil {
			t.Fatal("Unable to store:", err)
		}

		user2 := entity.User{Id: 2, Name: "Subi"}
		post2 := entity.Post{
			User: user2,
			Date: "2022-19-02",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if err = repo.StoreLike(post); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		if err = repo.StoreDislike(post2); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		if liked, err := repo.FetchIdsByReaction(user, "liked"); err != nil {
			t.Fatal("Unable to FetchIdsByReaction:", err)
		} else if liked[0] != 1 {
			t.Fatalf("want id = %d, got id = %d:", post.Id, liked[0])
		}

		if disliked, err := repo.FetchIdsByReaction(user2, "disliked"); err != nil {
			t.Fatal("Unable to FetchIdsByReaction:", err)
		} else if disliked[0] != 2 {
			t.Fatalf("want id = %d, got id = %d:", post2.Id, disliked[0])
		}
	})
}

func TestPostGetById(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		title := "Travel"
		post1 := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Title:   title,
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		title2 := "Sports"
		post2 := entity.Post{
			User:    entity.User{Id: 2, Name: "Subi"},
			Date:    "2022-19-02",
			Title:   title2,
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Title != title {
			t.Fatalf("want title = %v, got title = %v:", title, found.Title)
		}

		if found, err := repo.GetById(2); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Title != title2 {
			t.Fatalf("want title = %v, got title = %v:", title2, found.Title)
		}
	})
}

func TestGetRelatedCategories(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		categories := []string{"Cars", "Cinema", "Games"}
		post1 := entity.Post{
			User:       entity.User{Id: 5, Name: "Riddle"},
			Categories: categories,
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}
		if err := repo.StoreTopicReference(post1); err != nil {
			t.Fatal("Unable to StoreTopicReference:", err)
		}

		categories2 := []string{"Games", "Work"}
		post2 := entity.Post{
			User:       entity.User{Id: 5, Name: "Riddle"},
			Categories: categories2,
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		}
		if err := repo.StoreTopicReference(post2); err != nil {
			t.Fatal("Unable to StoreTopicReference:", err)
		}

		if found, err := repo.GetRelatedCategories(entity.Post{Id: 1}); err != nil {
			t.Fatal("Unable to GetRelatedCategories:", err)
		} else if !reflect.DeepEqual(found, categories) {
			t.Fatalf("want = %v, got = %v:", categories, found)
		}

		if found, err := repo.GetRelatedCategories(entity.Post{Id: 2}); err != nil {
			t.Fatal("Unable to GetRelatedCategories:", err)
		} else if !reflect.DeepEqual(found, categories2) {
			t.Fatalf("want = %v, got = %v:", categories2, found)
		}
	})
}

func TestPostUpdate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		title := "Travel"
		post1 := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Title:   title,
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		title2 := "Sports"
		post2 := entity.Post{
			User:    entity.User{Id: 2, Name: "Subi"},
			Date:    "2022-19-02",
			Title:   title2,
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post2); err != nil {
			t.Fatal("Unable to store:", err)
		}

		newTitle := "NewTitle"
		if err = repo.Update(entity.Post{Id: 1, Title: newTitle}); err != nil {
			t.Fatal("Unable to Update:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Title != newTitle {
			t.Fatalf("want title = %v, got title = %v:", newTitle, found.Title)
		}
	})
}

func TestPostDelete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		post1 := entity.Post{
			User:    entity.User{Id: 1, Name: "Riddle"},
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.Id != 1 {
			t.Fatalf("want id = %d, got id = %d:", found.Id, 1)
		}

		if err = repo.Delete(entity.Post{Id: 1}); err != nil {
			t.Fatal("Unable to Delete:", err)
		}

		expErr := "no rows in result set"

		if _, err := repo.GetById(1); err == nil {
			t.Fatal("Expected error:", err)
		} else if !strings.Contains(err.Error(), expErr) {
			t.Fatalf("want err = %v, got err = %v:", expErr, err)
		}
	})
}

func TestPostDeleteLikes(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		user1 := entity.User{Id: 1, Name: "Riddle"}
		user2 := entity.User{Id: 2, Name: "Subi"}
		post1 := entity.Post{
			User:    user1,
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if err = repo.StoreLike(post1); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		post1.User = user2

		if err = repo.StoreLike(post1); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.TotalLikes != 2 {
			t.Fatalf("want len = %d, got len = %d:", 2, found.TotalLikes)
		}

		if err = repo.DeleteLike(post1); err != nil {
			t.Fatal("Unable to DeleteLike:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.TotalLikes != 1 {
			t.Fatalf("want len = %d, got len = %d:", 1, found.TotalLikes)
		}
	})
}

func TestPostDeleteDislikes(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		user1 := entity.User{Id: 1, Name: "Riddle"}
		user2 := entity.User{Id: 2, Name: "Subi"}
		post1 := entity.Post{
			User:    user1,
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if err = repo.StoreDislike(post1); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		post1.User = user2

		if err = repo.StoreDislike(post1); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.TotalDislikes != 2 {
			t.Fatalf("want len = %d, got len = %d:", 2, found.TotalDislikes)
		}

		if err = repo.DeleteDislike(post1); err != nil {
			t.Fatal("Unable to DeleteDislike:", err)
		}

		if found, err := repo.GetById(1); err != nil {
			t.Fatal("Unable to GetById:", err)
		} else if found.TotalDislikes != 1 {
			t.Fatalf("want len = %d, got len = %d:", 1, found.TotalDislikes)
		}
	})
}

func TestPostFetchReactions(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		user := entity.User{Id: 5, Name: "Ids"}
		user1 := entity.User{Id: 1, Name: "Riddle"}
		user2 := entity.User{Id: 2, Name: "Subi"}
		post1 := entity.Post{
			User:    user1,
			Date:    "2022-19-01",
			Content: "Lorem ipsum dolor sit amet.",
		}

		if err := repo.Store(&post1); err != nil {
			t.Fatal("Unable to store:", err)
		}

		if err = repo.StoreDislike(post1); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		post1.User = user2

		if err = repo.StoreDislike(post1); err != nil {
			t.Fatal("Unable to StoreDislike:", err)
		}

		post1.User = user

		if err = repo.StoreLike(post1); err != nil {
			t.Fatal("Unable to StoreLike:", err)
		}

		if found, err := repo.FetchReactions(1); err != nil {
			t.Fatal("Unable to FetchReactions:", err)
		} else if len(found.Dislikes) != 2 {
			t.Fatalf("want dislikes = %d, got dislikes = %d:", 2, len(found.Dislikes))
		} else if len(found.Likes) != 1 {
			t.Fatalf("want likes = %d, got likes = %d:", 1, len(found.Likes))
		}
	})
}

func TestStoreCategories(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
		defer sqlite.MustCloseDB(t, db)
		err := sqlite.CreateDB(db)
		if err != nil {
			t.Fatal("Unable to create db:", err)
		}
		repo := sqlite.NewPostsRepo(db)

		categories1 := []string{"Cars", "Sports"}
		categories2 := []string{"Cars", "Travel"}
		summary := []string{"Cars", "Sports", "Travel"}

		if err = repo.StoreCategories(categories1); err != nil {
			t.Fatal("Unable to StoreCategories:", err)
		}

		if found, err := repo.GetExistedCategories(); err != nil {
			t.Fatal("Unable to GetExistedCategories:", err)
		} else if !reflect.DeepEqual(found, categories1) {
			t.Fatalf("want categories = %v, got categories = %v:", categories1, found)
		}

		if err = repo.StoreCategories(categories2); err != nil {
			t.Fatal("Unable to StoreCategories:", err)
		}

		if found, err := repo.GetExistedCategories(); err != nil {
			t.Fatal("Unable to GetExistedCategories:", err)
		} else if !reflect.DeepEqual(found, summary) {
			t.Fatalf("want categories = %v, got categories = %v:", summary, found)
		}
	})
}
