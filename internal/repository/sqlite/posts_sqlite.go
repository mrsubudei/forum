package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/entity"
	"forum/pkg/sqlite3"
	"time"
)

type PostsRepo struct {
	*sqlite3.Sqlite
}

func NewPostsRepo(sq *sqlite3.Sqlite) *PostsRepo {
	return &PostsRepo{sq}
}

func (pr *PostsRepo) Store(post *entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - Begin: %w", err)
	}

	stmt, err := tx.Prepare(
		`INSERT INTO posts(user_id, date, title, content) 
			values(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()
	date := post.Date.Format(TimeFormat)
	res, err := stmt.Exec(post.User.Id, date, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("PostsRepo - Store - RowsAffected: %w", err)
	}
	postId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - LastInsertId: %w", err)
	}
	post.Id = postId
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) StoreTopicReference(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreTopicReference - Begin: %w", err)
	}

	stmt, err := tx.Prepare(
		`INSERT INTO reference_topic(post_id, topic) 
			values(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < len(post.Category); i++ {
		res, err := stmt.Exec(post.Id, post.Category[i])
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreTopicReference - Exec: %w", err)
		}
		affected, err := res.RowsAffected()
		if affected != 1 || err != nil {
			return fmt.Errorf("PostsRepo - StoreTopicReference - RowsAffected: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreTopicReference - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) Fetch() ([]entity.Post, error) {
	var posts []entity.Post
	return posts, nil
}

//todo
func (pr *PostsRepo) GetById(id int64) (entity.Post, error) {
	var post entity.Post

	stmt, err := pr.DB.Prepare(`
	SELECT
	user_id, date, title, content,
	(SELECT name FROM users WHERE users.id = user_id) AS user_name,
	(SELECT date FROM post_likes WHERE post_likes.post_id = ?) AS post_likes,
	(SELECT date FROM post_dislikes WHERE post_likes.post_id = ?) AS post_dislikes,
	FROM posts
	WHERE id = ?
	`)

	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()
	var postsLikes sql.NullInt64
	var postDislikes sql.NullInt64
	var date string

	err = stmt.QueryRow(id, id, id).Scan(&post.User.Id, &date, &post.Title, &post.Content,
		&post.User.Name, &post.Likes, &post.Dislikes)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Scan: %w", err)
	}
	dateParsed, err := time.Parse(DateParseFormat, date)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Parse regDate: %w", err)
	}

	post.TotalLikes = postsLikes.Int64
	post.TotalDislikes = postDislikes.Int64
	post.Date = dateParsed

	return post, nil
}

func (pr *PostsRepo) GetIdByCategory(category string) (entity.Post, error) {
	var post entity.Post

	stmt, err := pr.DB.Prepare(`
	SELECT
	post_id
	FROM reference_topic
	WHERE topic = ?
	`)

	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(category).Scan(&post.Id)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Scan: %w", err)
	}

	return post, nil
}

func (pr *PostsRepo) GetRelatedCategories(post entity.Post) ([]string, error) {
	categories := []string{}
	rows, err := pr.DB.Query(`
	SELECT
	topic
	FROM reference_topic
	WHERE post_id = ?
	`)

	if err != nil {
		return nil, fmt.Errorf("PostsRepo - GetRelatedCategories - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err = rows.Scan(category)
		if err != nil {
			return nil, fmt.Errorf("PostsRepo - GetRelatedCategories - Scan: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (pr *PostsRepo) Update(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	UPDATE posts
	SET title = ?, content = ?
	WHERE id = ?
	`)

	if err != nil {
		return fmt.Errorf("UsersRepo - Update - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Title, post.Content, post.Id)
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

func (pr *PostsRepo) Delete(post entity.Post) error {
	return nil
}

func (pr *PostsRepo) StoreLike(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreLike - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO post_likes(post_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	date := post.Date.Format(TimeFormat)

	res, err := stmt.Exec(post.Id, post.User.Id, date)
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreLike - Exec: %w", err)
	}
	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("PostsRepo - StoreLike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreLike - Commit: %w", err)
	}
	return nil
}

func (pr *PostsRepo) DeleteLike(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	DELETE FROM post_likes
	WHERE post_id = ?
	`)

	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) StoreDislike(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreDislike - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO post_dislikes(post_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	date := post.Date.Format(TimeFormat)

	res, err := stmt.Exec(post.Id, post.User.Id, date)
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreDislike - Exec: %w", err)
	}
	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("PostsRepo - StoreDislike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreDislike - Commit: %w", err)
	}
	return nil
}

func (pr *PostsRepo) DeleteDislike(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	DELETE FROM post_dislikes
	WHERE post_id = ?
	`)

	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Commit: %w", err)
	}

	return nil
}
