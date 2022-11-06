package sqlite

import (
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

	stmt, err := tx.Prepare(
		`INSERT INTO posts(user_id, date, title, content) 
			values(?, ?, ?, ?)`)
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

	stmt, err := tx.Prepare(
		`INSERT INTO reference_topic(post_id, topic_id) 
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

func (pr *PostsRepo) FetchLiked() ([]entity.Post, error) {
	var posts []entity.Post
	return posts, nil
}

func (pr *PostsRepo) FetchDisliked() ([]entity.Post, error) {
	var posts []entity.Post
	return posts, nil
}

func (pr *PostsRepo) GetById(n int) (entity.Post, error) {
	var post entity.Post
	return post, nil
}

func (pr *PostsRepo) GetByCategory(category string) (entity.Post, error) {
	var post entity.Post
	return post, nil
}

func (pr *PostsRepo) Update(post entity.Post) (entity.Post, error) {
	return post, nil
}

func (pr *PostsRepo) Delete(post entity.Post) error {
	return nil
}

func (pr *PostsRepo) ThumbsUp(post entity.Post) error {
	return nil
}

func (pr *PostsRepo) ThumbsDown(post entity.Post) error {
	return nil
}
