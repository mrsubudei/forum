package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/entity"
	"forum/pkg/sqlite3"
	"time"
)

type UsersRepo struct {
	*sqlite3.Sqlite
}

func NewUsersRepo(sq *sqlite3.Sqlite) *UsersRepo {
	return &UsersRepo{sq}
}

func (ur *UsersRepo) Store(user entity.User) error {
	tx, err := ur.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO users(name, email, password, reg_date, date_of_birth, city, sex, role, sign) 
		values(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		return fmt.Errorf("UsersRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.RegDate,
		user.DateOfBirth, user.City, user.Gender, RoleUser, " ")
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

	rows, err := ur.DB.Query(`
	SELECT
		id, name, email, reg_date, date_of_birth, city, sex, role,
		(SELECT id FROM posts WHERE posts.user_id = users.id) AS posts,
		(SELECT id FROM comments WHERE comments.user_id = users.id) AS comments,
		(SELECT post_id FROM post_likes WHERE post_likes.user_id = users.id) AS post_likes,
		(SELECT post_id FROM post_dislikes WHERE post_dislikes.user_id = users.id) AS post_dislikes,
		(SELECT comment_id FROM comment_likes WHERE comment_likes.user_id = users.id) AS comment_likes,
		(SELECT comment_id FROM comment_dislikes WHERE comment_dislikes.user_id = users.id) AS comment_dislikes
	FROM users
	`)
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
			&user.Gender, &user.Role, &posts, &comments, &postLikes, &postDislikes, &commentLikes, &commentDislikes)
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

func (ur *UsersRepo) GetId(user entity.User) (int64, error) {
	var id int64
	switch {
	case user.SessionToken != "":
		stmt, err := ur.DB.Prepare(`
		SELECT id
		FROM users
		WHERE session_token = ?
		`)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Session - Query: %w", err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(user.SessionToken).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Session - Scan: %w", err)
		}
	case user.Name != "":
		stmt, err := ur.DB.Prepare(`
		SELECT id
		FROM users
		WHERE name = ?
		`)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Name - Query: %w", err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(user.Name).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Name - Scan: %w", err)
		}
	case user.Email != "":
		stmt, err := ur.DB.Prepare(`
		SELECT id
		FROM users
		WHERE email = ?
		`)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Email - Query: %w", err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(user.Email).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("UsersRepo - GetId - case Email - Scan: %w", err)
		}
	}

	return id, nil
}

func (ur *UsersRepo) GetById(id int64) (entity.User, error) {
	var user entity.User
	stmt, err := ur.DB.Prepare(`
	SELECT
		id, name, email, password, reg_date, date_of_birth, city, sex, role, sign,
		(SELECT COUNT(*) FROM posts WHERE posts.user_id = users.id) AS posts,
		(SELECT COUNT(*) FROM comments WHERE comments.user_id = users.id) AS comments,
		(SELECT COUNT(*) FROM post_likes WHERE post_likes.user_id = users.id) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes WHERE post_dislikes.user_id = users.id) AS post_dislikes,
		(SELECT COUNT(*) FROM comment_likes WHERE comment_likes.user_id = users.id) AS comment_likes,
		(SELECT COUNT(*) FROM comment_dislikes WHERE comment_dislikes.user_id = users.id) AS comment_dislikes
	FROM users
	WHERE id = ?
	`)
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
	var sign sql.NullString

	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate, &user.DateOfBirth, &user.City,
		&user.Gender, &user.Role, &sign, &posts, &comments, &postLikes, &postDislikes, &commentLikes, &commentDislikes)
	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetById - Scan: %w", err)
	}

	user.Posts = posts.Int64
	user.Comments = comments.Int64
	user.PostLikes = postLikes.Int64
	user.PostDislikes = postDislikes.Int64
	user.CommentLikes = commentLikes.Int64
	user.CommentDislikes = commentDislikes.Int64
	user.Sign = sign.String

	if user.DateOfBirth == "0001-01-01" {
		user.DateOfBirth = ""
	}

	return user, nil
}

func (ur *UsersRepo) GetSession(n int64) (entity.User, error) {
	var user entity.User
	stmt, err := ur.DB.Prepare(`
	SELECT
		session_token, session_ttl
	FROM users
	WHERE id = ?
	`)
	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetSession - Query: %w", err)
	}
	defer stmt.Close()
	var sessionToken sql.NullString
	var sessionTTL sql.NullString
	err = stmt.QueryRow(n).Scan(&sessionToken, &sessionTTL)
	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetSession - Scan: %w", err)
	}
	if sessionTTL.String == "" {
		return user, entity.ErrUserNotFound
	}
	TTLParsed, err := time.Parse(DateAndTimeFormat, sessionTTL.String)
	if err != nil {
		return user, fmt.Errorf("UsersRepo - GetSession - Parse TTL: %w", err)
	}
	user.SessionTTL = TTLParsed
	user.SessionToken = sessionToken.String
	return user, nil
}

func (ur *UsersRepo) UpdateInfo(user entity.User) error {
	tx, err := ur.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`
	UPDATE users
	SET date_of_birth = ?, city = ?, sex = ?, sign = ?, role = ?
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.DateOfBirth, user.City, user.Gender, user.Sign, user.Role, user.Id)
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
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdatePassword - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`
	UPDATE users
	SET password = ?
	WHERE id = ?
	`)
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

func (ur *UsersRepo) NewSession(user entity.User) error {
	tx, err := ur.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - NewSession - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`
	UPDATE users
	SET session_token = ?, session_ttl = ?
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("UsersRepo - NewSession - Prepare: %w", err)
	}
	defer stmt.Close()
	sessionTTL := user.SessionTTL.Format(DateAndTimeFormat)
	res, err := stmt.Exec(user.SessionToken, sessionTTL, user.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - NewSession - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - NewSession - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - NewSession - Commit: %w", err)
	}

	return nil
}

func (ur *UsersRepo) UpdateSession(user entity.User) error {
	tx, err := ur.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdateSession - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`
	UPDATE users
	SET session_ttl = ?
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdateSession - Prepare: %w", err)
	}
	defer stmt.Close()
	sessionTTL := user.SessionTTL.Format(DateAndTimeFormat)
	res, err := stmt.Exec(sessionTTL, user.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdateSession - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - UpdateSession - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - UpdateSession - Commit: %w", err)
	}

	return nil
}

func (ur *UsersRepo) Delete(user entity.User) error {
	tx, err := ur.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("UsersRepo - Delete - Begin: %w", err)
	}
	stmt, err := ur.DB.Prepare(`
	DELETE FROM users
	WHERE id = ?
	`)
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
