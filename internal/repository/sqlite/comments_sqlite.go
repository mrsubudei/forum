package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/entity"
	"forum/pkg/sqlite3"
)

type CommentsRepo struct {
	*sqlite3.Sqlite
}

func NewCommentsRepo(sq *sqlite3.Sqlite) *CommentsRepo {
	return &CommentsRepo{sq}
}

func (cr *CommentsRepo) Store(comment entity.Comment) error {
	tx, err := cr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Store - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO comments(post_id, user_id, date, content) 
		values(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.PostId, comment.User.Id, comment.Date, comment.Content)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Store - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentsRepo - Store - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Store - Commit: %w", err)
	}

	return nil
}

func (cr *CommentsRepo) Fetch(postId int64) ([]entity.Comment, error) {
	var comments []entity.Comment

	rows, err := cr.DB.Query(`
	SELECT
		id, post_id, user_id, date, content,
		(SELECT COUNT(*) FROM comment_likes WHERE comment_likes.comment_id = comments.id) AS comment_likes,
		(SELECT COUNT(*) FROM comment_dislikes WHERE comment_dislikes.comment_id = comments.id) AS comment_dislikes
	FROM comments
	WHERE post_id = ?
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("CommentsRepo - Fetch - Query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment entity.Comment
		var commentLikes sql.NullInt64
		var commentDislikes sql.NullInt64

		err = rows.Scan(&comment.Id, &comment.PostId, &comment.User.Id, &comment.Date, &comment.Content,
			&commentLikes, &commentDislikes)
		if err != nil {
			return nil, fmt.Errorf("CommentsRepo - Fetch - Scan: %w", err)
		}

		comment.TotalLikes = commentLikes.Int64
		comment.TotalDislikes = commentDislikes.Int64
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cr *CommentsRepo) GetById(commentId int64) (entity.Comment, error) {
	var comment entity.Comment

	stmt, err := cr.DB.Prepare(`
	SELECT
		id, post_id, user_id, date, content
		(SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?) AS comment_likes,
		(SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = ?) AS comment_dislikes
	FROM comments
	WHERE id = ?
	`)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()
	var commentLikes sql.NullInt64
	var commentDislikes sql.NullInt64
	err = stmt.QueryRow(commentId).Scan(&comment.Id, &comment.PostId, &comment.User.Id, &comment.Date, &comment.Content, &commentLikes, &commentDislikes)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - GetById - Scan: %w", err)
	}

	comment.TotalLikes = commentLikes.Int64
	comment.TotalDislikes = commentDislikes.Int64

	return comment, nil
}

func (cr *CommentsRepo) GetPostIds(user entity.User) ([]int64, error) {
	var postIds []int64

	rows, err := cr.DB.Query(`
	SELECT DISTINCT post_id
	FROM comments
	WHERE user_id = ?
	`, user.Id)
	if err != nil {
		return nil, fmt.Errorf("PostsRepo - GetPostIds - Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postId sql.NullInt64
		err = rows.Scan(&postId)
		if err != nil {
			return postIds, fmt.Errorf("PostsRepo - GetPostIds - Scan: %w", err)
		}
		postIds = append(postIds, postId.Int64)
	}

	return postIds, nil
}

func (cr *CommentsRepo) Update(comment entity.Comment) error {
	tx, err := cr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Update - Begin: %w", err)
	}
	stmt, err := cr.DB.Prepare(`
	UPDATE comments
	SET content = ?
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Update - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.Content, comment.Id)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Update - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentsRepo - Update - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Update - Commit: %w", err)
	}

	return nil
}

func (cr *CommentsRepo) Delete(comment entity.Comment) error {
	tx, err := cr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Delete - Begin: %w", err)
	}
	stmt, err := cr.DB.Prepare(`
	DELETE FROM comments
	WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Delete - Prepare: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.Id)
	if err != nil {
		return fmt.Errorf("CommentsRepo - Delete - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentsRepo - Delete - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - Delete - Commit: %w", err)
	}

	return nil
}

func (pr *CommentsRepo) StoreLike(comment entity.Comment) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - StoreLike - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO comment_likes(comment_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.Id, comment.User.Id, getRegTime(DateFormat))
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("CommentsRepo - StoreLike - Exec - Commit: %w", err)
		}
		return fmt.Errorf("CommentsRepo - StoreLike - Exec: %w", err)
	}
	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentsRepo - StoreLike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - StoreLike - Commit: %w", err)
	}
	return nil
}

func (pr *CommentsRepo) DeleteLike(comment entity.Comment) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteLike - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	DELETE FROM comment_likes
	WHERE comment_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteLike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Id, comment.User.Id)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteLike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteLike - Commit: %w", err)
	}

	return nil
}

func (pr *CommentsRepo) StoreDislike(comment entity.Comment) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - StoreDislike - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO comment_dislikes(comment_id, user_id, date) 
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(comment.Id, comment.User.Id, getRegTime(DateFormat))
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("CommentsRepo - StoreDislike - Exec - Commit: %w", err)
		}
		return fmt.Errorf("CommentsRepo - StoreDislike - Exec: %w", err)
	}
	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentsRepo - StoreDislike - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - StoreDislike - Commit: %w", err)
	}
	return nil
}

func (pr *CommentsRepo) DeleteDislike(comment entity.Comment) error {
	tx, err := pr.DB.Begin()
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Begin: %w", err)
	}
	stmt, err := pr.DB.Prepare(`
	DELETE FROM comment_dislikes
	WHERE comment_id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Id, comment.User.Id)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Commit: %w", err)
	}

	return nil
}

func (pr *CommentsRepo) FetchReactions(id int64) (entity.Comment, error) {
	var comment entity.Comment
	var likes []entity.Reaction
	var dislikes []entity.Reaction

	rowsLikes, err := pr.DB.Query(`
		SELECT user_id, date
		FROM comment_likes
		WHERE comment_id = ?
	`, id)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - FetchReactions - likes - Query: %w", err)
	}
	defer rowsLikes.Close()

	for rowsLikes.Next() {
		var like entity.Reaction
		err = rowsLikes.Scan(&like.UserId, &like.Date)
		if err != nil {
			return comment, fmt.Errorf("CommentsRepo - FetchReactions - likes - Scan: %w", err)
		}
		likes = append(likes, like)
	}

	rowsDislikes, err := pr.DB.Query(`
		SELECT user_id, date
		FROM comment_dislikes
		WHERE comment_id = ?
	`, id)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - FetchReactions - dislikes - Query: %w", err)
	}
	defer rowsDislikes.Close()

	for rowsDislikes.Next() {
		var dislike entity.Reaction
		err = rowsDislikes.Scan(&dislike.UserId, &dislike.Date)
		if err != nil {
			return comment, fmt.Errorf("CommentsRepo - FetchReactions - dislikes - Scan: %w", err)
		}
		dislikes = append(dislikes, dislike)
	}
	comment.Likes = append(comment.Likes, likes...)
	comment.Dislikes = append(comment.Dislikes, dislikes...)
	return comment, nil
}
