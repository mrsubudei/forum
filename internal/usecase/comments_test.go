package usecase_test

import (
	"testing"

	"forum/internal/entity"
	m "forum/internal/repository/sqlite/mock"
	"forum/internal/usecase"
)

var (
	comment1 = entity.Comment{
		PostId:  1,
		User:    user1,
		Date:    "2022-05-04",
		Content: "Lorem ipsum, dolor sit amet",
	}
	comment2 = entity.Comment{
		PostId:  1,
		User:    user4,
		Date:    "2022-05-04",
		Content: "Lorem sit amet",
	}
)

func TestWriteComments(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		commentUseCase := usecase.NewCommentsUseCase(mockRepo.Comments, mockRepo.Posts, mockRepo.Users)
		if err := commentUseCase.WriteComment(comment1); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAllComments(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)
		commentUseCase := usecase.NewCommentsUseCase(mockRepo.Comments, mockRepo.Posts, mockRepo.Users)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}

		if err := commentUseCase.WriteComment(comment1); err != nil {
			t.Fatal(err)
		}

		if err := commentUseCase.WriteComment(comment2); err != nil {
			t.Fatal(err)
		}

		if found, err := commentUseCase.GetAllComments(post1.Id); err != nil {
			t.Fatal(err)
		} else if len(found) != 2 {
			t.Fatalf("want: %d, got: %d", 2, len(found))
		}
	})
}

func TestUpdateComment(t *testing.T) {
	mockRepo := m.NewMockRepos()
	commentUseCase := usecase.NewCommentsUseCase(mockRepo.Comments, mockRepo.Posts, mockRepo.Users)

	t.Run("OK", func(t *testing.T) {
		if err := commentUseCase.WriteComment(comment1); err != nil {
			t.Fatal(err)
		}

		newContent := "New content abc"
		comment1.Content = newContent

		if err := commentUseCase.UpdateComment(comment1); err != nil {
			t.Fatal(err)
		}
		if mockRepo.Comments.Comments[0].Content != newContent {
			t.Fatalf("want: %v, got: %v", newContent, mockRepo.Comments.Comments[0].Content)
		}
	})

	t.Run("err not found", func(t *testing.T) {
		notExisted := entity.Comment{
			Id:      15,
			Content: "abv",
			User:    user1,
		}

		if err := commentUseCase.UpdateComment(notExisted); err == nil {
			t.Fatal("Epected error")
		}
	})
}

func TestDeleteComment(t *testing.T) {
	mockRepo := m.NewMockRepos()
	commentUseCase := usecase.NewCommentsUseCase(mockRepo.Comments, mockRepo.Posts, mockRepo.Users)

	t.Run("OK", func(t *testing.T) {
		if err := commentUseCase.WriteComment(comment1); err != nil {
			t.Fatal(err)
		}

		if len(mockRepo.Comments.Comments) != 1 {
			t.Fatalf("want: %d, got: %d", 1, len(mockRepo.Comments.Comments))
		}

		if err := commentUseCase.DeleteComment(comment1); err != nil {
			t.Fatal(err)
		}

		if len(mockRepo.Comments.Comments) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(mockRepo.Comments.Comments))
		}
	})

	t.Run("err not exist", func(t *testing.T) {
		if err := commentUseCase.DeleteComment(comment1); err == nil {
			t.Fatal("Expected error")
		}
	})
}

func TestCommentMakeReaction(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)
		commentUseCase := usecase.NewCommentsUseCase(mockRepo.Comments, mockRepo.Posts, mockRepo.Users)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user5); err != nil {
			t.Fatal(err)
		}

		if err := commentUseCase.WriteComment(comment1); err != nil {
			t.Fatal(err)
		}

		if err := commentUseCase.MakeReaction(comment1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		comment1.User = user4
		if err := commentUseCase.MakeReaction(comment1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		if found, err := commentUseCase.GetReactions(comment1.Id, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 2 {
			t.Fatalf("want: %d, got: %d", 2, len(found))
		}

		if err := commentUseCase.MakeReaction(comment1, usecase.ReactionLike); err != nil {
			t.Fatal(err)
		}

		if found, err := commentUseCase.GetReactions(comment1.Id, usecase.PostLikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 1 {
			t.Fatalf("want: %d, got: %d", 1, len(found))
		}

		comment1.User = user5

		if err := commentUseCase.MakeReaction(comment1, usecase.ReactionDislike); err != nil {
			t.Fatal(err)
		}

		if found, err := commentUseCase.GetReactions(comment1.Id, usecase.PostDislikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 1 {
			t.Fatalf("want: %d, got: %d", 1, len(found))
		}

		if err := commentUseCase.MakeReaction(comment1, usecase.ReactionDislike); err != nil {
			t.Fatal(err)
		}

		if found, err := commentUseCase.GetReactions(comment1.Id, usecase.PostDislikedQuery); err != nil {
			t.Fatal(err)
		} else if len(found) != 0 {
			t.Fatalf("want: %d, got: %d", 0, len(found))
		}
	})
}
