package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"strings"
	"time"

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

const (
	UpdateInfoQuery     = "info"
	UpdatePasswordQuery = "password"
	UpdateSessionQuery  = "session"
	UniqueEmailErr      = "UNIQUE constraint failed: users.email"
	UniqueNameErr       = "UNIQUE constraint failed: users.name"
	DateFormat          = "2006-01-02"
	DateAndTimeFormat   = "2006-01-02 15:04:05"
)

func NewUsersUseCase(repo repository.Users, hasher hasher.PasswordHasher,
	tokenManager auth.TokenManager, postUseCase repository.Posts,
	commentUseCase repository.Comments) *UsersUseCase {

	return &UsersUseCase{
		repo:           repo,
		hasher:         hasher,
		tokenManager:   tokenManager,
		postUseCase:    postUseCase,
		commentUseCase: commentUseCase,
	}
}

func (uu *UsersUseCase) SignUp(user entity.User) error {

	hashed, err := uu.hasher.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignUp #1 - %w", err)
	}
	user.Password = hashed

	user.RegDate = getRegTime(DateAndTimeFormat)

	err = uu.repo.Store(user)
	if err != nil {
		if strings.Contains(err.Error(), UniqueEmailErr) {
			return entity.ErrUserEmailAlreadyExists
		}
		if strings.Contains(err.Error(), UniqueNameErr) {
			return entity.ErrUserNameAlreadyExists
		}
		return fmt.Errorf("UsersUseCase - SignUp #2 - %w", err)
	}

	return nil
}

func (uu *UsersUseCase) SignIn(user entity.User) error {
	id, err := uu.repo.GetId(user)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn #1 - %w", err)
	}

	if id == 0 {
		return entity.ErrUserNotFound
	}

	existUserInfo, err := uu.GetById(id)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn #2 - %w", err)
	}

	err = uu.hasher.CheckPassword(existUserInfo.Password, user.Password)
	if err != nil {
		return entity.ErrUserPasswordIncorrect
	}
	token, err := uu.tokenManager.NewToken()
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn #3 - %w", err)
	}

	user.SessionToken = token
	TTL := uu.tokenManager.UpdateTTL()
	user.SessionTTL = TTL
	user.Id = id

	err = uu.repo.NewSession(user)
	if err != nil {
		return fmt.Errorf("UsersUseCase - SignIn #4 - %w", err)
	}
	return nil
}

func (uu *UsersUseCase) GetIdBy(user entity.User) (int64, error) {
	id, err := uu.repo.GetId(user)

	if err != nil {
		return 0, fmt.Errorf("UsersUseCase - GetIdBy - %w", err)
	}

	if id == 0 {
		return 0, entity.ErrUserNotFound
	}
	return id, nil
}

func (uu *UsersUseCase) UpdateSession(user entity.User) error {
	TTL := uu.tokenManager.UpdateTTL()
	user.SessionTTL = TTL
	err := uu.UpdateUserInfo(user, UpdateSessionQuery)
	if err != nil {
		return fmt.Errorf("UsersUseCase - UpdateSession - %w", err)
	}

	return nil
}

func (uu *UsersUseCase) DeleteSession(user entity.User) error {
	timeNow := time.Now()
	formatted := timeNow.Format(DateAndTimeFormat)
	parsed, err := time.Parse(DateAndTimeFormat, formatted)
	if err != nil {
		return fmt.Errorf("UsersUseCase - DeleteSession #1 - %w", err)
	}
	user.SessionTTL = parsed
	err = uu.UpdateUserInfo(user, UpdateSessionQuery)
	if err != nil {
		return fmt.Errorf("UsersUseCase - DeleteSession #2 - %w", err)
	}

	return nil
}

func (uu *UsersUseCase) GetSession(id int64) (entity.User, error) {
	var user entity.User
	user, err := uu.repo.GetSession(id)
	if err != nil {
		if strings.Contains(err.Error(), NoRowsResultErr) {
			return user, entity.ErrUserNotFound
		}
		return user, fmt.Errorf("UsersUseCase - GetSession - %w", err)
	}
	user.Id = id
	return user, nil
}

func (uu *UsersUseCase) CheckSession(user entity.User) (bool, error) {
	if user.Id == 0 {
		return false, nil
	}
	existUserInfo, err := uu.GetSession(user.Id)
	if err != nil {
		return false, fmt.Errorf("UsersUseCase - CheckSession #1 - %w", err)
	}

	//check token
	if existUserInfo.SessionToken != user.SessionToken {
		return false, nil
	}

	//check token life time
	expired, err := uu.tokenManager.CheckTTLExpired(existUserInfo.SessionTTL)
	if err != nil {
		if strings.Contains(err.Error(), NoRowsResultErr) {
			return false, entity.ErrUserNotFound
		}
		return false, fmt.Errorf("UsersUseCase - CheckSession #2 - %w", err)
	}

	return !expired, nil
}

func (uu *UsersUseCase) GetAllUsers() ([]entity.User, error) {
	users, err := uu.repo.Fetch()
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - GetAllUsers - %w", err)
	}

	return users, nil
}

func (uu *UsersUseCase) GetById(id int64) (entity.User, error) {
	var user entity.User
	user, err := uu.repo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), NoRowsResultErr) {
			return user, entity.ErrUserNotFound
		}
		return user, fmt.Errorf("UsersUseCase - GetById - %w", err)
	}
	user.Id = id
	return user, nil
}

func (uu *UsersUseCase) UpdateUserInfo(user entity.User, query string) error {
	switch query {
	case UpdateInfoQuery:
		err := uu.repo.UpdateInfo(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo #1 - %w", err)
		}
	case UpdatePasswordQuery:
		hashed, err := uu.hasher.Hash(user.Password)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo #2 - %w", err)
		}
		user.Password = hashed
		err = uu.repo.UpdatePassword(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo #3 - %w", err)
		}
	case UpdateSessionQuery:
		err := uu.repo.UpdateSession(user)
		if err != nil {
			return fmt.Errorf("UsersUseCase - UpdateUserInfo #4 - %w", err)
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

func getRegTime(format string) string {
	timeNow := time.Now()
	return timeNow.Format(format)
}
