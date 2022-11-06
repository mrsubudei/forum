package sqlite

import (
	"database/sql"
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

func (ur *UsersRepo) Store(user entity.User) error {
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

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.RegDate,
		user.DateOfBirth, user.City, user.Sex)
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - Store - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Commit: %w", err)
	}

	return nil
}

func (ur *UsersRepo) Fetch() ([]entity.User, error) {
	var users []entity.User

	rows, err := ur.DB.Query(`SELECT
			id, name, email, reg_date, date_of_birth, city, sex,
			(SELECT id FROM posts WHERE posts.user_id = users.id) AS posts,
			(SELECT id FROM comments WHERE comments.user_id = users.id) AS comments,
			(SELECT post_id FROM post_likes WHERE post_likes.user_id = users.id) AS post_likes,
			(SELECT post_id FROM post_dislikes WHERE post_dislikes.user_id = users.id) AS post_dislikes,
			(SELECT comment_id FROM comment_likes WHERE comment_likes.user_id = users.id) AS comment_likes,
			(SELECT comment_id FROM comment_dislikes WHERE comment_dislikes.user_id = users.id) AS comment_dislikes
			FROM users`)

	if err != nil {
		return nil, fmt.Errorf("UsersRepo - Fetch - Query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := entity.User{}
		var posts sql.NullInt64
		var comments sql.NullInt64
		var postLikes sql.NullInt64
		var postDislikes sql.NullInt64
		var commentLikes sql.NullInt64
		var commentDislikes sql.NullInt64
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.RegDate, &user.DateOfBirth, &user.City,
			&user.Sex, &posts, &comments, &postLikes, &postDislikes, &commentLikes, &commentDislikes)
		if err != nil {
			return nil, fmt.Errorf("UsersRepo - Fetch - Scan: %w", err)
		}

		user.Posts = posts.Int64
		user.Comments = comments.Int64
		user.PostLikes = postLikes.Int64
		user.PostDislikes = postDislikes.Int64
		user.CommentLikes = commentLikes.Int64
		user.CommentDislikes = commentDislikes.Int64
		users = append(users, user)
	}

	return users, nil
}

func (ur *UsersRepo) GetById(n int) (entity.User, error) {
	var user entity.User
	stmt, err := ur.DB.Prepare(`SELECT
	name, email, reg_date, date_of_birth, city, sex,
	(SELECT id FROM posts WHERE posts.user_id = users.id) AS posts,
	(SELECT id FROM comments WHERE comments.user_id = users.id) AS comments,
	(SELECT post_id FROM post_likes WHERE post_likes.user_id = users.id) AS post_likes,
	(SELECT post_id FROM post_dislikes WHERE post_dislikes.user_id = users.id) AS post_dislikes,
	(SELECT comment_id FROM comment_likes WHERE comment_likes.user_id = users.id) AS comment_likes,
	(SELECT comment_id FROM comment_dislikes WHERE comment_dislikes.user_id = users.id) AS comment_dislikes
	FROM users
	WHERE id = ?`)

	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()
	var posts sql.NullInt64
	var comments sql.NullInt64
	var postLikes sql.NullInt64
	var postDislikes sql.NullInt64
	var commentLikes sql.NullInt64
	var commentDislikes sql.NullInt64
	err = stmt.QueryRow(n).Scan(&user.Name, &user.Email, &user.RegDate, &user.DateOfBirth, &user.City,
		&user.Sex, &posts, &comments, &postLikes, &postDislikes, &commentLikes, &commentDislikes)

	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetById - Scan: %w", err)
	}
	user.Posts = posts.Int64
	user.Comments = comments.Int64
	user.PostLikes = postLikes.Int64
	user.PostDislikes = postDislikes.Int64
	user.CommentLikes = commentLikes.Int64
	user.CommentDislikes = commentDislikes.Int64

	return user, nil
}

func (ur *UsersRepo) UpdateInfo(user entity.User) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`UPDATE users
	SET email = ?, date_of_birth = ?, city = ?, sex = ?
	WHERE id = ?`)

	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Email, user.DateOfBirth, user.City, user.Sex, user.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - Update - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Commit: %w", err)
	}
	return nil
}

func (ur *UsersRepo) UpdatePassword(user entity.User) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`UPDATE users
	SET password = ?
	WHERE id = ?`)

	if err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Password, user.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - Commit: %w", err)
	}
	return nil
}

func (ur *UsersRepo) Delete(user entity.User) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - Delete - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`DELETE FROM users
	WHERE id = ?`)

	if err != nil {
		return fmt.Errorf("UsersRepo - Delete - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - Delete - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - Delete - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - Delete - Commit: %w", err)
	}
	return nil
}
