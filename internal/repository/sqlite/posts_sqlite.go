package sqlite

import (
	"forum/internal/entity"
	"forum/pkg/sqlite3"
)

type PostsRepo struct {
	*sqlite3.Sqlite
}

func NewPostsRepo(sq *sqlite3.Sqlite) *PostsRepo {
	return &PostsRepo{sq}
}

func (pr *PostsRepo) Store(p entity.Post) error {
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
