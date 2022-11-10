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

	stmt, err := tx.Prepare(`
	INSERT INTO posts(user_id, date, title, content) 
		values(?, ?, ?, ?)
	`)
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

	stmt, err := tx.Prepare(`
	INSERT INTO reference_topic(post_id, topic) 
		values(?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < len(post.Categories); i++ {
		res, err := stmt.Exec(post.Id, post.Categories[i])
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

	rows, err := pr.DB.Query(`
	SELECT
		id, user_id, date, title, content,
		(SELECT name FROM users WHERE users.id = posts.user_id) AS user_name,
		(SELECT COUNT(*) FROM post_likes ) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes ) AS post_dislikes
	FROM posts
	`)

	if err != nil {
		return nil, fmt.Errorf("PostsRepo - Fetch - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		var postsLikes sql.NullInt64
		var postDislikes sql.NullInt64
		var date string
		err = rows.Scan(&post.Id, &post.User.Id, &date, &post.Title, &post.Content,
			&post.User.Name, &postsLikes, &postDislikes)
		if err != nil {
			return posts, fmt.Errorf("PostsRepo - Fetch - Scan: %w", err)
		}
		dateParsed, err := time.Parse(DateParseFormat, date)
		if err != nil {
			return posts, fmt.Errorf("PostsRepo - Fetch - Parse date: %w", err)
		}
		post.TotalLikes = postsLikes.Int64
		post.TotalDislikes = postDislikes.Int64
		post.Date = dateParsed
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostsRepo) GetById(id int64) (entity.Post, error) {
	var post entity.Post

	stmt, err := pr.DB.Prepare(`
	SELECT
		user_id, date, title, content,
		(SELECT name FROM users WHERE users.id = posts.user_id) AS user_name,
		(SELECT COUNT(*) FROM post_likes ) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes ) AS post_dislikes
	FROM posts
	WHERE id = ?
	`)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Prepare: %w", err)
	}
	defer stmt.Close()
	var postsLikes sql.NullInt64
	var postDislikes sql.NullInt64
	var date string

	err = stmt.QueryRow(id).Scan(&post.User.Id, &date, &post.Title, &post.Content,
		&post.User.Name, &postsLikes, &postDislikes)
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

func (pr *PostsRepo) GetIdByCategory(category string) (int64, error) {
	var id int64

	stmt, err := pr.DB.Prepare(`
	SELECT post_id
	FROM reference_topic
	WHERE topic = ?
	`)
	if err != nil {
		return 0, fmt.Errorf("PostsRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(category).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("PostsRepo - GetById - Scan: %w", err)
	}

	return id, nil
}

func (pr *PostsRepo) GetRelatedCategories(post entity.Post) ([]string, error) {

	categories := []string{}
	rows, err := pr.DB.Query(`
	SELECT topic
	FROM reference_topic
	WHERE post_id = ?
	`, post.Id)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - GetRelatedCategories - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err = rows.Scan(&category)
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
		return fmt.Errorf("PostsRepo - Update - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	UPDATE posts
	SET title = ?, content = ?
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - Update - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Title, post.Content, post.Id)
	if err != nil {
		return fmt.Errorf("PostsRepo - Update - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("PostsRepo - Update - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - Update - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) Delete(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - Delete - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	DELETE FROM posts
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - Delete - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Id)
	if err != nil {
		return fmt.Errorf("PostsRepo - Delete - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("PostsRepo - Delete - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - Delete - Commit: %w", err)
	}

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
		tx.Commit()
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreLike - Exec err Commit: %w", err)
		}
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
	WHERE post_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.User.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteLike - Exec: %w", err)
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
		tx.Commit()
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreDislike - Exec err Commit: %w", err)
		}
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
	WHERE post_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.User.Id)
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UsersRepo - DeleteDislike - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) FetchReactions(id int64) (entity.Post, error) {
	var post entity.Post
	var likes []entity.Reaction
	var dislikes []entity.Reaction

	rowsLikes, err := pr.DB.Query(`
		SELECT user_id, date
		FROM post_likes
		WHERE post_id = ?
	`, id)

	if err != nil {
		return post, fmt.Errorf("PostsRepo - FetchReactions - likes - Query: %w", err)
	}
	defer rowsLikes.Close()

	for rowsLikes.Next() {
		var like entity.Reaction
		var date string
		err = rowsLikes.Scan(&like.UserId, &date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - likes - Scan: %w", err)
		}
		dateParsed, err := time.Parse(DateParseFormat, date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - likes - Parse date: %w", err)
		}
		like.Date = dateParsed
		likes = append(likes, like)
	}

	rowsDislikes, err := pr.DB.Query(`
		SELECT user_id, date
		FROM post_dislikes
		WHERE post_id = ?
	`, id)

	if err != nil {
		return post, fmt.Errorf("PostsRepo - FetchReactions - dislikes - Query: %w", err)
	}
	defer rowsDislikes.Close()

	for rowsDislikes.Next() {
		var dislike entity.Reaction
		var date string
		err = rowsDislikes.Scan(&dislike.UserId, &date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - dislikes - Scan: %w", err)
		}
		dateParsed, err := time.Parse(DateParseFormat, date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - dislikes - Parse date: %w", err)
		}
		dislike.Date = dateParsed
		likes = append(likes, dislike)
	}
	post.Likes = append(post.Likes, likes...)
	post.Dislikes = append(post.Dislikes, dislikes...)
	return post, nil
}
