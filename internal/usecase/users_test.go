package usecase_test

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"forum/internal/config"
	"forum/internal/entity"
	m "forum/internal/repository/sqlite/mock"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
)

var errNoRows = "no rows in result set"

var (
	user1 = entity.User{
		Id:       1,
		Name:     "Riddle",
		Email:    "Riddle@mail.ru",
		Password: "Vivse",
		City:     "Astana",
		Role:     "Пользователь",
	}
	user2 = entity.User{
		Id:       2,
		Name:     "Riddle",
		Email:    "Subi@mail.ru",
		City:     "Karaganda",
		Password: "azh",
		Role:     "Пользователь",
	}
	user3 = entity.User{
		Id:       3,
		Name:     "Isa",
		Email:    "Subi@mail.ru",
		Password: "Mimi",
		City:     "Karaganda",
		Role:     "Пользователь",
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
	t.Parallel()
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err name already exist", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

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
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

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

func TestSignIn(t *testing.T) {
	t.Parallel()
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(entity.User{Name: "Riddle", Password: "Vivse"}); err != nil {
			t.Fatal(err)
		}
		if err := userUseCase.SignIn(entity.User{Email: "Riddle@mail.ru", Password: "Vivse"}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err user not exist", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager, mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}
		if err := userUseCase.SignIn(user3); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want err: %v, got: %v", entity.ErrUserNotFound, err)
		}
	})

	t.Run("err wrong password", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(entity.User{Name: "Riddle", Password: "abc"}); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserPasswordIncorrect) {
			t.Fatalf("want err: %v, got: %v", entity.ErrUserPasswordIncorrect, err)
		}
	})
}

func TestUsersGetIdBy(t *testing.T) {
	t.Parallel()
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if id, err := userUseCase.GetIdBy(entity.User{Name: "Riddle"}); err != nil {
			t.Fatal(err)
		} else if id != user1.Id {
			t.Fatalf("want id: %d, got: %d", user1.Id, id)
		}

		if id, err := userUseCase.GetIdBy(entity.User{Email: "Riddle@mail.ru"}); err != nil {
			t.Fatal(err)
		} else if id != user1.Id {
			t.Fatalf("want id: %d, got: %d", user1.Id, id)
		}

		if err := userUseCase.SignIn(entity.User{Name: "Riddle", Password: "Vivse"}); err != nil {
			t.Fatal(err)
		}

		token := mockRepo.Users.Users[0].SessionToken

		if id, err := userUseCase.GetIdBy(entity.User{SessionToken: token}); err != nil {
			t.Fatal(err)
		} else if id != user1.Id {
			t.Fatalf("want id: %d, got: %d", user1.Id, id)
		}
	})

	t.Run("err user not found", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if _, err := userUseCase.GetIdBy(entity.User{Name: "Pup"}); err == nil {
			t.Fatal("Expected error")
		} else if !strings.Contains(err.Error(), errNoRows) {
			t.Fatalf("want err: %v, got: %v", errNoRows, err)
		}
	})
}

func TestUpdateSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		TTL := mockRepo.Users.Users[0].SessionTTL

		if err := userUseCase.UpdateSession(user1); err != nil {
			t.Fatal(err)
		} else if !mockRepo.Users.Users[0].SessionTTL.After(TTL) {
			t.Fatal("could not update TTL")
		}
	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		TTL := mockRepo.Users.Users[0].SessionTTL.UTC().Format(usecase.DateAndTimeFormat)
		fmt.Println(TTL)
		if err := userUseCase.DeleteSession(user1); err != nil {
			t.Fatal(err)
		} else if mockRepo.Users.Users[0].SessionTTL.Format(usecase.DateAndTimeFormat) <= TTL {
			t.Fatal("Could not delete session")
		}
	})
}

func TestGetSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		token := mockRepo.Users.Users[0].SessionToken

		if found, err := userUseCase.GetSession(user1.Id); err != nil {
			t.Fatal(err)
		} else if found.SessionToken != token {
			t.Fatalf("want : %v, got: %v", token, found.SessionToken)
		}
	})
}

func TestCheckSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		if auth, err := userUseCase.CheckSession(user1); err != nil {
			t.Fatal(err)
		} else if auth {
			t.Fatalf("expected true")
		}
	})

	t.Run("err did not find session", func(t *testing.T) {
		hasher, tokenManager := getDependencies()
		mockRepo := m.NewMockRepos()
		userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
			mockRepo.Posts, mockRepo.Comments)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		if _, err := userUseCase.CheckSession(user2); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want : %v, got: %v", entity.ErrUserNotFound, err)
		}
	})
}
