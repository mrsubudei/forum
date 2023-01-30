package usecase_test

import (
	"errors"
	"log"
	"testing"
	"time"

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
		Email:    "Riddle@mail.ru",
		Password: "Mimi",
		City:     "Karaganda",
		Role:     "Пользователь",
	}

	user4 = entity.User{
		Id:       4,
		Name:     "Makaron",
		Email:    "Makaron@mail.ru",
		Password: "Mimi",
		Role:     "Пользователь",
	}

	user5 = entity.User{
		Id:       5,
		Name:     "Spaget",
		Email:    "Spaget@mail.ru",
		Password: "Mimi",
		Role:     "Пользователь",
	}
)

func setupUserUseCase(mockRepo *m.MockRepos) *usecase.UsersUseCase {
	cfg, err := config.LoadConfig("../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	hasher := hasher.NewBcryptHasher()
	tokenManager := auth.NewManager(cfg)

	userUseCase := usecase.NewUsersUseCase(mockRepo.Users, hasher, tokenManager,
		mockRepo.Posts, mockRepo.Comments)
	return userUseCase
}

func TestSignUp(t *testing.T) {
	mockRepo := m.NewMockRepos()
	userUseCase := setupUserUseCase(mockRepo)

	t.Run("OK", func(t *testing.T) {
		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err name already exist", func(t *testing.T) {
		if err := userUseCase.SignUp(user2); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNameAlreadyExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrUserNameAlreadyExists, err)
		}
	})

	t.Run("err email already exist", func(t *testing.T) {
		if err := userUseCase.SignUp(user3); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserEmailAlreadyExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrUserEmailAlreadyExists, err)
		}
	})
}

func TestSignIn(t *testing.T) {
	mockRepo := m.NewMockRepos()
	userUseCase := setupUserUseCase(mockRepo)

	t.Run("OK", func(t *testing.T) {
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
		if err := userUseCase.SignIn(user3); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want err: %v, got: %v", entity.ErrUserNotFound, err)
		}
	})

	t.Run("err wrong password", func(t *testing.T) {
		if err := userUseCase.SignIn(entity.User{Name: "Riddle", Password: "abc"}); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserPasswordIncorrect) {
			t.Fatalf("want err: %v, got: %v", entity.ErrUserPasswordIncorrect, err)
		}
	})
}

func TestUsersGetIdBy(t *testing.T) {
	mockRepo := m.NewMockRepos()
	userUseCase := setupUserUseCase(mockRepo)

	t.Run("OK", func(t *testing.T) {
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
		if _, err := userUseCase.GetIdBy(entity.User{Name: "Pup"}); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want err: %v, got: %v", errNoRows, err)
		}
	})
}

func TestUpdateSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		TTL := mockRepo.Users.Users[0].SessionTTL.Add(-time.Hour * 10)

		if err := userUseCase.UpdateSession(user1); err != nil {
			t.Fatal(err)
		} else if !mockRepo.Users.Users[0].SessionTTL.After(TTL) {
			t.Fatal("could not update TTL")
		}
	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		TTL := mockRepo.Users.Users[0].SessionTTL
		if err := userUseCase.DeleteSession(user1); err != nil {
			t.Fatal(err)
		} else if !TTL.After(mockRepo.Users.Users[0].SessionTTL) {
			t.Fatal("Could not delete session")
		}
	})
}

func TestGetSession(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

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
	mockRepo := m.NewMockRepos()
	userUseCase := setupUserUseCase(mockRepo)

	t.Run("OK", func(t *testing.T) {
		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		if auth, err := userUseCase.CheckSession(mockRepo.Users.Users[0]); err != nil {
			t.Fatal(err)
		} else if !auth {
			t.Fatal("expected true")
		}
	})

	t.Run("err did not find session", func(t *testing.T) {
		if _, err := userUseCase.CheckSession(user2); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want : %v, got: %v", entity.ErrUserNotFound, err)
		}
	})

	t.Run("err expired session", func(t *testing.T) {
		mockRepo.Users.Users[0].SessionTTL = mockRepo.Users.Users[0].SessionTTL.Add(-time.Hour * 25)

		if auth, err := userUseCase.CheckSession(mockRepo.Users.Users[0]); err != nil {
			t.Fatal(err)
		} else if auth {
			t.Fatal("expected false")
		}
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}
		if err := userUseCase.SignUp(user5); err != nil {
			t.Fatal(err)
		}

		if users, err := userUseCase.GetAllUsers(); err != nil {
			t.Fatal(err)
		} else if len(users) != 3 {
			t.Fatalf("want: %d, got: %d", 3, len(users))
		}
	})
}

func TestUserGetById(t *testing.T) {
	mockRepo := m.NewMockRepos()
	userUseCase := setupUserUseCase(mockRepo)

	t.Run("OK", func(t *testing.T) {
		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignUp(user4); err != nil {
			t.Fatal(err)
		}

		if found, err := userUseCase.GetById(4); err != nil {
			t.Fatal(err)
		} else if found.Id != 4 {
			t.Fatalf("want: %d, got: %d", 4, found.Id)
		}
	})

	t.Run("err not found", func(t *testing.T) {
		if _, err := userUseCase.GetById(3); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("want: %v, got: %v", entity.ErrUserNotFound, err)
		}
	})
}

func TestUpdateUserInfo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		city := "Updated"
		user1.City = city
		if err := userUseCase.UpdateUserInfo(user1, "info"); err != nil {
			t.Fatal(err)
		}
		if found, err := userUseCase.GetById(1); err != nil {
			t.Fatal(err)
		} else if found.City != city {
			t.Fatalf("want: %v, got: %v", city, found.City)
		}

		password := "NewPassword"
		user1.Password = password
		if err := userUseCase.UpdateUserInfo(user1, "password"); err != nil {
			t.Fatal(err)
		}
		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal("Could not update password")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockRepo := m.NewMockRepos()
		userUseCase := setupUserUseCase(mockRepo)

		if err := userUseCase.SignUp(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.DeleteUser(user1); err != nil {
			t.Fatal(err)
		}

		if err := userUseCase.SignIn(user1); err == nil {
			t.Fatal("Could not delete user")
		}
	})
}
