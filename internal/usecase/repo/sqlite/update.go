package sqlite

import (
	"forum/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

func (c *CommunicationRepo) CreateUser(u *entity.User) (int, error) {
	tx, err := c.DB.Begin()
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

	cached, err := cachePass(u.Password)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(u.Name, u.Email, cached, u.RegDate, u.DateOfBirth, u.City, u.Sex)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userId), nil
}

func (c *CommunicationRepo) CreatePost(p *entity.Post) (int, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO posts(user_id, date, content) 
			values(?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.Date, p.Content)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postId), nil
}

func (c *CommunicationRepo) CreateComment(p *entity.Post, date, content string) error {
	tx, err := c.DB.Begin()
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

func (c *CommunicationRepo) PutPostLike(u *entity.User, postId int, date string) error {
	tx, err := c.DB.Begin()
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

func (c *CommunicationRepo) PutPostDisLike(u *entity.User, postId int, date string) error {
	tx, err := c.DB.Begin()
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

func (c *CommunicationRepo) PutCommentLike(u *entity.User, commentId int, date string) error {
	tx, err := c.DB.Begin()
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

func (c *CommunicationRepo) PutCommentDisLike(u *entity.User, commentId int, date string) error {
	tx, err := c.DB.Begin()
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

func (c *CommunicationRepo) CreateTopics(name []string) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	for i := 0; i < len(name); i++ {
		stmt, err := tx.Prepare(
			`INSERT INTO topics(name) 
				values(?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(name[i])
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

func (c *CommunicationRepo) CreatePostRef(p *entity.Post, name []string) error {
	tx, err := c.DB.Begin()
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

func cachePass(str string) (string, error) {
	cached, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(cached), nil
}
