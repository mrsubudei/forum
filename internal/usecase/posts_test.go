package usecase_test

import (
	"errors"
	"testing"

	"forum/internal/entity"
	m "forum/internal/repository/sqlite/mock"
	"forum/internal/usecase"
)

var (
	post1 = entity.Post{
		Id:         1,
		User:       user1,
		Categories: []string{"Cars"},
		Date:       "2022-05-02",
		Title:      "Audi",
		Content:    "Lorem ipsum, dolor sit amet consectetur adipisicing.",
	}

	post2 = entity.Post{
		Id:         2,
		User:       user4,
		Date:       "2022-01-03",
		Categories: []string{"Sports"},
		Title:      "Footbal",
		Content:    "Lorem ipsum, dolor sit amet.",
	}

	post3 = entity.Post{
		Id:         3,
		User:       user1,
		Date:       "2022-04-02",
		Title:      "BMW",
		Categories: []string{"Cars"},
		Content:    "Lorem ipsum, dolor sit.",
	}

	post4 = entity.Post{
		Id:         4,
		User:       user1,
		Date:       "2022-04-02",
		Title:      "Bdrg",
		Categories: []string{"Guns", "Computers"},
		Content:    "Lorem ipsum, dolor sit.",
	}
)

func TestCreatePost(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAllPost(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post2); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetAllPosts(); err != nil {
			t.Fatal(err)
		} else if len(found) != 2 {
			t.Fatalf("want: %d, got: %d", 2, len(found))
		}
	})
}

func TestGetPostsByQuery(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post2); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post3); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostAuthorQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 2 {
			t.Fatalf("want: %d, got: %d", 2, len(found))
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostDislikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostCommentedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}
	})
}

func TestPostGetById(t *testing.T) {
	mockRepo := m.NewMockRepos()
	hasher, tokenManager := getDependencies()
	postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
	userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
		mockRepo.Posts, mockRepo.Comments)

	t.Run("OK", func(t *testing.T) {
		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post2); err != nil {
			t.Fatal(err)
		}

		id := post2.Id
		if found, err := postUseCase.GetById(id); err != nil {
			t.Fatal(err)
		} else if found.Id != id {
			t.Fatalf("want: %d, got: %d", id, found.Id)
		}
	})

	t.Run("err not found", func(t *testing.T) {
		id := int64(14)
		if _, err := postUseCase.GetById(id); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrPostNotFound) {
			t.Fatalf("want: %v, got: %v", entity.ErrPostNotFound, err)
		}
	})
}

func TestGetAllByCategory(t *testing.T) {
	mockRepo := m.NewMockRepos()
	hasher, tokenManager := getDependencies()
	postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
	userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
		mockRepo.Posts, mockRepo.Comments)

	t.Run("OK", func(t *testing.T) {
		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post2); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post3); err != nil {
			t.Fatal(err)
		}

		category1 := post1.Categories[0]
		category2 := post2.Categories[0]

		if found, err := postUseCase.GetAllByCategory(category1); err != nil {
			t.Fatal(err)
		} else if len(found) != 2 {
			t.Fatalf("want: %d, got: %d", 2, len(found))
		}

		if found, err := postUseCase.GetAllByCategory(category2); err != nil {
			t.Fatal(err)
		} else if found[0].Categories[0] != category2 {
			t.Fatalf("want: %v, got: %v", category2, found[0].Categories[0])
		}
	})

	t.Run("err not found", func(t *testing.T) {
		category := "Weather"

		if _, err := postUseCase.GetAllByCategory(category); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrPostNotFound) {
			t.Fatalf("want: %v, got: %v", entity.ErrPostNotFound, err)
		}
	})
}

func TestGetAllCategories(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		mockRepo.Posts.AllTopics = map[string]bool{}
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)

		categories := []string{}
		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := mockRepo.Posts.StoreCategories(post1.Categories); err != nil {
			t.Fatal(err)
		}
		categories = append(categories, post1.Categories...)

		if err := postUseCase.CreatePost(post4); err != nil {
			t.Fatal(err)
		}

		if err := mockRepo.Posts.StoreCategories(post4.Categories); err != nil {
			t.Fatal(err)
		}
		categories = append(categories, post4.Categories...)

		if found, err := postUseCase.GetAllCategories(); err != nil {
			t.Fatal(err)
		} else if len(found) != len(categories) {
			t.Fatalf("want: %d, got: %d", len(categories), len(found))
		}
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		hasher, tokenManager := getDependencies()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		content := "New content"

		if err := postUseCase.UpdatePost(entity.Post{Id: 1, Content: content, User: user1}); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetById(1); err != nil {
			t.Fatal(err)
		} else if found.Content != content {
			t.Fatal("Could not update")
		}
	})
}

func TestDeletePost(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		hasher, tokenManager := getDependencies()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if _, err := postUseCase.GetById(1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.DeletePost(post1); err != nil {
			t.Fatal(err)
		}

		if _, err := postUseCase.GetById(1); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrPostNotFound) {
			t.Fatalf("want: %v, got: %v", entity.ErrPostNotFound, err)
		}
	})
}

func TestPostMakeReaction(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		hasher, tokenManager := getDependencies()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post2); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.MakeReaction(post1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if found[0].Id != post1.Id {
			t.Fatalf("want: %d, got: %d", post1.Id, found[0].Id)
		}

		if err := postUseCase.MakeReaction(post2, usecase.ReactionDislike); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user4, usecase.PostDislikedQuery); err != nil {
			t.Fatal(err)
		} else if found[0].Id != post2.Id {
			t.Fatalf("want: %d, got: %d", post2.Id, found[0].Id)
		}

		// second time reaction should delete reaction
		if err := postUseCase.MakeReaction(post2, usecase.ReactionDislike); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user4, usecase.PostDislikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}
	})
}

func TestPostDeleteReaction(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		hasher, tokenManager := getDependencies()
		postUseCase := usecase.NewPostsUseCase(mockRepo.Posts, mockRepo.Users, mockRepo.Comments)
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.CreatePost(post1); err != nil {
			t.Fatal(err)
		}

		if err := postUseCase.MakeReaction(post1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 1 {
			t.Fatalf("want: %d, got: %d", 1, len(found))
		}

		if err := postUseCase.DeleteReaction(post1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		if found, err := postUseCase.GetPostsByQuery(user1, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}
	})
}
