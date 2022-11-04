package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository/sqlite"
	"forum/pkg/auth"
	"forum/pkg/hasher"
)

type UsersUseCase struct {
	repo           sqlite.Users
	hasher         hasher.PasswordHasher
	tokenManager   auth.TokenManager
	postUseCase    sqlite.Posts
	commentUseCase sqlite.Comments
}

func NewUsersUseCase(repo sqlite.Users, hasher hasher.PasswordHasher,
	tokenManager auth.TokenManager, postUseCase sqlite.Posts, commentUseCase sqlite.Comments) *UsersUseCase {
	return &UsersUseCase{
		repo:           repo,
		hasher:         hasher,
		tokenManager:   tokenManager,
		postUseCase:    postUseCase,
		commentUseCase: commentUseCase,
	}
}

func (uu *UsersUseCase) SignUp(u entity.User) error {
	err := uu.repo.Store(u)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignUp - %w", err)
	}
	return nil
}

func (uu *UsersUseCase) SignIn(u entity.User) error {
	return nil
}

func (uu *UsersUseCase) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	return users, nil
}

func (uu *UsersUseCase) GetOne() (entity.User, error) {
	var user entity.User
	return user, nil
}

func (uu *UsersUseCase) UpdateUserInfo() (entity.User, error) {
	var user entity.User
	return user, nil
}

func (uu *UsersUseCase) DeleteUser(u entity.User) error {
	return nil
}