package usecase_test

import (
	"errors"
	"log"
	"testing"

	"forum/internal/config"
	"forum/internal/entity"
	m "forum/internal/repository/sqlite/mock"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
)

var (
	user1 = entity.User{
		Id:    1,
		Name:  "Riddle",
		Email: "Riddle@mail.ru",
		City:  "Astana",
		Role:  "Пользователь",
		Sign:  " ",
	}
	user2 = entity.User{
		Id:    2,
		Name:  "Riddle",
		Email: "Subi@mail.ru",
		City:  "Karaganda",
		Role:  "Пользователь",
		Sign:  " ",
	}
	user3 = entity.User{
		Id:    3,
		Name:  "Isa",
		Email: "Subi@mail.ru",
		City:  "Karaganda",
		Role:  "Пользователь",
		Sign:  " ",
	}
)

func getDependencies() (*hasher.BcryptHasher, *auth.Manager) {
	cfg, err := config.LoadConfig("../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	hasher := hasher.NewBcryptHasher()
	tokenManager := auth.NewManager(cfg)
	return hasher, tokenManager
}

func TestSignUp(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		usersRepo := m.NewUsersMockRepo()
		postsRepo := m.NewPostsMockrepo()
		commentsRepo := m.NewCommentsMockrepo()

		userUseCase := usecase.NewUsersUseCase(usersRepo, hasher, tokenManager, postsRepo, commentsRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err name already exist", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mr := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mr.Users, hasher, tokenManager, mr.Posts, mr.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user2); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNameAlreadyExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrUserNameAlreadyExists, err)
		}
	})

	t.Run("err email already exist", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mr := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mr.Users, hasher, tokenManager, mr.Posts, mr.Comments)

		if err := userUseCase.SignUp(user2); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user3); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserEmailAlreadyExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrUserEmailAlreadyExists, err)
		}
	})
}
