package sqlite

import (
	"database/sql"
	"fmt"

	"forum/internal/entity"
	"forum/pkg/sqlite3"
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
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO posts(user_id, date, title, content) 
		values(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.User.Id, post.Date, post.Title, post.Content)
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
	defer tx.Rollback()

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
		(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes WHERE post_dislikes.post_id = posts.id) AS post_dislikes
	FROM posts
	`)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - Fetch - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		var userName sql.NullString
		var postsLikes sql.NullInt64
		var postDislikes sql.NullInt64

		err = rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content,
			&userName, &postsLikes, &postDislikes)
		if err != nil {
			return posts, fmt.Errorf("PostsRepo - Fetch - Scan: %w", err)
		}

		post.TotalLikes = postsLikes.Int64
		post.TotalDislikes = postDislikes.Int64
		post.User.Name = userName.String

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostsRepo) FetchByAuthor(user entity.User) ([]entity.Post, error) {
	var posts []entity.Post

	rows, err := pr.DB.Query(`
	SELECT
		id, user_id, date, title, content,
		(SELECT name FROM users WHERE users.id = posts.user_id) AS user_name,
		(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes WHERE post_dislikes.post_id = posts.id) AS post_dislikes
	FROM posts
	WHERE user_id = ?
	`, user.Id)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - FetchByQuery - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		var userName sql.NullString
		var postsLikes sql.NullInt64
		var postDislikes sql.NullInt64

		err = rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content,
			&userName, &postsLikes, &postDislikes)
		if err != nil {
			return posts, fmt.Errorf("PostsRepo - FetchByQuery - Scan: %w", err)
		}

		post.TotalLikes = postsLikes.Int64
		post.TotalDislikes = postDislikes.Int64
		post.User.Name = userName.String

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostsRepo) FetchIdsByReaction(user entity.User, reaction string) ([]int64, error) {
	var postIds []int64
	var rows *sql.Rows
	var err error

	switch reaction {
	case QueryLiked:
		rows, err = pr.DB.Query(`
		SELECT post_id
		FROM post_likes
		WHERE user_id = ?
	`, user.Id)
	case QueryDislike:
		rows, err = pr.DB.Query(`
		SELECT post_id
		FROM post_dislikes
		WHERE user_id = ?
	`, user.Id)
	}

	if err != nil {
		return nil, fmt.Errorf("PostsRepo - FetchIdsByReaction - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postId sql.NullInt64

		err = rows.Scan(&postId)
		if err != nil {
			return postIds, fmt.Errorf("PostsRepo - FetchIdsByReaction - Scan: %w", err)
		}
		postIds = append(postIds, postId.Int64)
	}

	return postIds, nil
}

func (pr *PostsRepo) GetById(id int64) (entity.Post, error) {
	var post entity.Post

	stmt, err := pr.DB.Prepare(`
	SELECT
		id, user_id, date, title, content,
		(SELECT name FROM users WHERE users.id = posts.user_id) AS user_name,
		(SELECT COUNT(*) FROM post_likes WHERE post_id = ?) AS post_likes,
		(SELECT COUNT(*) FROM post_dislikes WHERE post_id = ?) AS post_dislikes
	FROM posts
	WHERE id = ?
	`)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Prepare: %w", err)
	}
	defer stmt.Close()
	var postsLikes sql.NullInt64
	var postDislikes sql.NullInt64
	var userName sql.NullString

	err = stmt.QueryRow(id, id, id).Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content,
		&userName, &postsLikes, &postDislikes)
	if err != nil {
		return post, fmt.Errorf("PostsRepo - GetById - Scan: %w", err)
	}

	post.TotalLikes = postsLikes.Int64
	post.TotalDislikes = postDislikes.Int64
	post.User.Name = userName.String

	return post, nil
}

func (pr *PostsRepo) GetIdsByCategory(category string) ([]int64, error) {
	var ids []int64

	rows, err := pr.DB.Query(`
	SELECT post_id
	FROM reference_topic
	WHERE topic = ?
	`, category)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - GetIdsByCategory - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("PostsRepo - GetIdsByCategory - Scan: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
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
	defer tx.Rollback()

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
	defer tx.Rollback()

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
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO post_likes(post_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Id, post.User.Id, getRegTime(DateFormat))
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreLike - Exec Commit: %w", err)
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
		return fmt.Errorf("PostsRepo - DeleteLike - Begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := pr.DB.Prepare(`
	DELETE FROM post_likes
	WHERE post_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteLike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.User.Id)
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteLike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteLike - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) StoreDislike(post entity.Post) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreDislike - Begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO post_dislikes(post_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Id, post.User.Id, getRegTime(DateFormat))
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreDislike - Exec Commit: %w", err)
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
		return fmt.Errorf("PostsRepo - DeleteDislike - Begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := pr.DB.Prepare(`
	DELETE FROM post_dislikes
	WHERE post_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteDislike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.User.Id)
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteDislike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - DeleteDislike - Commit: %w", err)
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
		err = rowsLikes.Scan(&like.UserId, &like.Date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - likes - Scan: %w", err)
		}
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
		err = rowsDislikes.Scan(&dislike.UserId, &dislike.Date)
		if err != nil {
			return post, fmt.Errorf("PostsRepo - FetchReactions - dislikes - Scan: %w", err)
		}
		dislikes = append(dislikes, dislike)
	}
	post.Likes = append(post.Likes, likes...)
	post.Dislikes = append(post.Dislikes, dislikes...)
	return post, nil
}

func (pr *PostsRepo) StoreCategories(categories []string) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreCategories - Begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO topics(name) 
		values(?)
	`)
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreCategories - Prepare: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < len(categories); i++ {
		res, err := stmt.Exec(categories[i])
		if err != nil {
			return fmt.Errorf("PostsRepo - StoreCategories - Exec: %w", err)
		}

		affected, err := res.RowsAffected()
		if affected != 1 || err != nil {
			return fmt.Errorf("PostsRepo - StoreCategories - RowsAffected: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PostsRepo - StoreCategories - Commit: %w", err)
	}

	return nil
}

func (pr *PostsRepo) GetExistedCategories() ([]string, error) {
	categories := []string{}
	rows, err := pr.DB.Query(`
	SELECT name
	FROM topics
	`)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - GetExistedCategories - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("PostsRepo - GetExistedCategories - Scan: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}
