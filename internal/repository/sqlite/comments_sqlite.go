package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/entity"
	"forum/pkg/sqlite3"
	"time"
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
		return fmt.Errorf("CommentssRepo - Store - Begin: %w", err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO comments(post_id, user_id, date, content) 
		values(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("CommentssRepo - Store - Prepare: %w", err)
	}
	defer stmt.Close()

	date := comment.Date.Format(TimeFormat)

	res, err := stmt.Exec(comment.Post.Id, comment.User.Id, date, comment.Content)
	if err != nil {
		return fmt.Errorf("CommentssRepo - Store - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("CommentssRepo - Store - RowsAffected: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentssRepo - Store - Commit: %w", err)
	}

	return nil
}

func (cr *CommentsRepo) Fetch() ([]entity.Comment, error) {
	var commets []entity.Comment

	rows, err := cr.DB.Query(`
	SELECT
		id, post_id, user_id, date, content,
		(SELECT COUNT(*) FROM comment_likes ) AS comment_likes,
		(SELECT COUNT(*) FROM comment_dislikes ) AS comment_dislikes
	FROM comments
	`)
	if err != nil {
		return nil, fmt.Errorf("CommentsRepo - Fetch - Query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment entity.Comment
		var commentLikes sql.NullInt64
		var commentDislikes sql.NullInt64
		var date string

		err = rows.Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &date, &comment.Content, &commentLikes, &commentDislikes)
		if err != nil {
			return nil, fmt.Errorf("CommentsRepo - Fetch - Scan: %w", err)
		}

		regDateParsed, err := time.Parse(DateParseFormat, date)
		if err != nil {
			return nil, fmt.Errorf("CommentsRepo - Fetch - Parse regDate: %w", err)
		}

		comment.TotalLikes = commentLikes.Int64
		comment.TotalDislikes = commentDislikes.Int64
		comment.Date = regDateParsed
		commets = append(commets, comment)
	}
	return commets, nil
}

func (cr *CommentsRepo) GetById(id int64) (entity.Comment, error) {
	var comment entity.Comment

	stmt, err := cr.DB.Prepare(`
	SELECT
		id, post_id, user_id, date, content
		(SELECT COUNT(*) FROM comment_likes WHERE post_likes.post_id = comments.post_id) AS comment_likes,
		(SELECT COUNT(*) FROM comment_dislikes WHERE post_likes.post_id = comments.post_id) AS comment_dislikes
	FROM comments
	WHERE id = ?
	`)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - GetById - Query: %w", err)
	}
	defer stmt.Close()
	var commentLikes sql.NullInt64
	var commentDislikes sql.NullInt64
	var date string
	err = stmt.QueryRow(id).Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &date, &comment.Content, &commentLikes, &commentDislikes)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - GetById - Scan: %w", err)
	}
	regDateParsed, err := time.Parse(DateParseFormat, date)
	if err != nil {
		return comment, fmt.Errorf("CommentsRepo - GetById - Parse regDate: %w", err)
	}

	comment.TotalLikes = commentLikes.Int64
	comment.TotalDislikes = commentDislikes.Int64
	comment.Date = regDateParsed

	return comment, nil
}

func (cr *CommentsRepo) Update(comment entity.Comment) (entity.Comment, error) {
	return comment, nil
}

func (cr *CommentsRepo) Delete(comment entity.Comment) error {
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

	date := comment.Date.Format(TimeFormat)

	res, err := stmt.Exec(comment.Id, comment.User.Id, date)
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("CommentsRepo - StoreLike - Exec err Commit: %w", err)
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
	WHERE comment_id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteLike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Id)
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

	date := comment.Date.Format(TimeFormat)

	res, err := stmt.Exec(comment.Id, comment.User.Id, date)
	if err != nil {
		tx.Commit()
		if err != nil {
			return fmt.Errorf("CommentsRepo - StoreDislike - Exec err Commit: %w", err)
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
	WHERE comment_id = ?
	`)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Id)
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CommentsRepo - DeleteDislike - Commit: %w", err)
	}

	return nil
}
