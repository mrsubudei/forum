package sqlite

import (
	"fmt"
	"forum/internal/entity"
	"forum/pkg/sqlite3"
)

type UsersRepo struct {
	*sqlite3.Sqlite
}

func NewUsersRepo(sq *sqlite3.Sqlite) *UsersRepo {
	return &UsersRepo{sq}
}

func (ur *UsersRepo) Store(u entity.User) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Begin: %w", err)
	}

	stmt, err := tx.Prepare(
		`INSERT INTO users(name, email, password, reg_date, date_of_birth, city, sex) 
			values(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Email, u.Password, u.RegDate, u.DateOfBirth, u.City, u.Sex)
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Commit: %w", err)
	}

	return nil
}

func (ur *UsersRepo) Fetch() ([]entity.User, error) {
	var users []entity.User
	return users, nil
}

func (ur *UsersRepo) GetById(n int) (entity.User, error) {
	var user entity.User
	return user, nil
}

func (ur *UsersRepo) Update(user entity.User) (entity.User, error) {
	return user, nil
}

func (ur *UsersRepo) Delete(user entity.User) error {
	return nil
}
