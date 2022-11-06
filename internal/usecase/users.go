package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"

	"forum/pkg/auth"
	"forum/pkg/hasher"
)

type UsersUseCase struct {
	repo           repository.Users
	hasher         hasher.PasswordHasher
	tokenManager   auth.TokenManager
	postUseCase    repository.Posts
	commentUseCase repository.Comments
}

func NewUsersUseCase(repo repository.Users, hasher hasher.PasswordHasher,
	tokenManager auth.TokenManager, postUseCase repository.Posts, commentUseCase repository.Comments) *UsersUseCase {
	return &UsersUseCase{
		repo:           repo,
		hasher:         hasher,
		tokenManager:   tokenManager,
		postUseCase:    postUseCase,
		commentUseCase: commentUseCase,
	}
}

func (uu *UsersUseCase) SignUp(u entity.User) error {
	hashed, err := uu.hasher.Hash(u.Password)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignUp - %w", err)
	}
	u.Password = hashed

	err = uu.repo.Store(u)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignUp - %w", err)
	}
	return nil
}

func (uu *UsersUseCase) SignIn(user entity.User) error {
	token, err := uu.tokenManager.NewToken()
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn - %w", err)
	}
	user.SessionToken = token
	TTL := uu.tokenManager.UpdateTTL()
	user.SessionTTL = TTL
	err = uu.repo.NewSession(user)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn - %w", err)
	}
	return nil
}

func (uu *UsersUseCase) UpdateSession(user entity.User) error {
	TTL := uu.tokenManager.UpdateTTL()
	user.SessionTTL = TTL
	err := uu.UpdateUserInfo(user, "session")
	if err != nil {
		return fmt.Errorf("UsersUseCase - UpdateSession - %w", err)
	}
	return nil
}

func (uu *UsersUseCase) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	users, err := uu.repo.Fetch()
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - GetAllUsers - %w", err)
	}
	return users, nil
}

func (uu *UsersUseCase) GetById(id int) (entity.User, error) {
	var user entity.User
	user, err := uu.repo.GetById(id)
	if err != nil {
		return user, fmt.Errorf("UsersUseCase - GetById - %w", err)
	}
	return user, nil
}

func (uu *UsersUseCase) GetSession(id int) (entity.User, error) {
	var user entity.User
	user, err := uu.repo.GetSession(id)
	if err != nil {
		return user, fmt.Errorf("UsersUseCase - GetSession - %w", err)
	}
	return user, nil
}

func (uu *UsersUseCase) UpdateUserInfo(user entity.User, query string) error {
	switch query {
	case "info":
		err := uu.repo.UpdateInfo(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo - %w", err)
		}
	case "password":
		hashed, err := uu.hasher.Hash(user.Password)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo - %w", err)
		}
		user.Password = hashed
		err = uu.repo.UpdatePassword(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo - %w", err)
		}
	case "session":
		err := uu.repo.UpdateSession(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo - %w", err)
		}
	}

	return nil
}

func (uu *UsersUseCase) DeleteUser(u entity.User) error {
	err := uu.repo.Delete(u)
	if err != nil {
		return fmt.Errorf("UsersUseCase - DeleteUser - %w", err)
	}
	return nil
}
