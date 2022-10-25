package usecases

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) CreateUser() (int, error) {

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return 0, err
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO users(name, email, password, reg_date, date_of_birth, city, sex) 
			values(?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	cached := cachePass(u.Password)

	_, err = stmt.Exec(u.Name, u.Email, cached, u.RegDate, u.DateOfBirth, u.City, u.Sex)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	//returning assigned id for user name
	row, err := db.Query("SELECT id FROM users ORDER BY ID DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var userId int
	if row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}

	return userId, nil
}

func (p *Post) CreatePost() (int, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return 0, err
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO posts(user_id, date, content) 
			values(?, ?, ?)`)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.UserId, p.Date, p.Content)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	//returning assigned id for post
	row, err := db.Query("SELECT id FROM posts ORDER BY ID DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var postId int
	if row.Next() {
		err = row.Scan(&postId)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}

	return postId, nil
}

func (p *Post) CreateComment(date, content string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO comments(post_id, user_id, date, content) 
			values(?, ?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Id, p.UserId, date, content)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PutPostLike(postId int, date string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO post_likes(post_id, user_id, date) 
			values(?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, u.Id, date)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PutPostDisLike(postId int, date string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO post_dislikes(post_id, user_id, date) 
			values(?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, u.Id, date)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PutCommentLike(commentId int, date string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO comment_likes(comment_id, user_id, date) 
			values(?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(commentId, u.Id, date)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PutCommentDisLike(commentId int, date string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO comment_dislikes(comment_id, user_id, date) 
			values(?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(commentId, u.Id, date)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func CreateTopics(name string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO topics(name) 
			values(?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) CreatePostRef(name ...string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO referencetopic(post_id, topic_id) 
			values(?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < len(name); i++ {
		_, err = stmt.Exec(p.Id, name[i])
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func cachePass(str string) string {
	cached, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(cached)
}
