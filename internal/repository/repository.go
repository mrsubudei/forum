package repository

import (
	"forum/internal/entity"
	"forum/internal/repository/sqlite"
	"forum/pkg/sqlite3"
)

type Posts interface {
	Store(p entity.Post) error
	Fetch() ([]entity.Post, error)
	FetchLiked() ([]entity.Post, error)
	FetchDisliked() ([]entity.Post, error)
	GetById(n int) (entity.Post, error)
	GetByCategory(category string) (entity.Post, error)
	Update(post entity.Post) (entity.Post, error)
	Delete(post entity.Post) error
	ThumbsUp(post entity.Post) error
	ThumbsDown(post entity.Post) error
}

type Users interface {
	Store(u entity.User) error
	Fetch() ([]entity.User, error)
	GetById(n int) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Delete(user entity.User) error
}

type Comments interface {
	Store(c entity.Comment) error
	Fetch() ([]entity.Comment, error)
	GetById(n int) (entity.Comment, error)
	Update(post entity.Comment) (entity.Comment, error)
	Delete(post entity.Comment) error
	ThumbsUp(comment entity.Comment) error
	ThumbsDown(comment entity.Comment) error
}

type UserFilter struct {
	Id                    []int
	Name                  string
	Email                 string
	RegDate               []string
	DateOfBirth           []string
	City                  []string
	Sex                   []string
	CountPosts            []int
	CountPostReactions    []int
	CountCommentReactions []int
}

type PostsFilter struct {
	Id            []int
	User          []int
	Date          []string
	Title         string
	Content       string
	Category      []string
	CountComments []int
	CountLikes    []int
	CountDislikes []int
}

type Repositories struct {
	Posts    Posts
	Users    Users
	Comments Comments
}

func NewRepositories(sq *sqlite3.Sqlite) *Repositories {
	return &Repositories{
		Posts:    sqlite.NewPostsRepo(sq),
		Users:    sqlite.NewUsersRepo(sq),
		Comments: sqlite.NewCommentsRepo(sq),
	}
}
