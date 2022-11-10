package repository

import (
	"forum/internal/entity"
	"forum/internal/repository/sqlite"
	"forum/pkg/sqlite3"
)

type Posts interface {
	Store(post *entity.Post) error
	Fetch() ([]entity.Post, error)
	GetById(id int64) (entity.Post, error)
	GetIdsByCategory(category string) ([]int64, error)
	Update(post entity.Post) error
	Delete(post entity.Post) error
	StoreLike(post entity.Post) error
	StoreDislike(post entity.Post) error
	DeleteLike(post entity.Post) error
	DeleteDislike(post entity.Post) error
	StoreTopicReference(post entity.Post) error
	GetRelatedCategories(post entity.Post) ([]string, error)
	FetchReactions(id int64) (entity.Post, error)
}

type Users interface {
	Store(user entity.User) error
	Fetch() ([]entity.User, error)
	GetId(user entity.User) (int64, error)
	GetById(n int64) (entity.User, error)
	GetSession(n int64) (entity.User, error)
	UpdateInfo(user entity.User) error
	UpdatePassword(user entity.User) error
	NewSession(user entity.User) error
	UpdateSession(user entity.User) error
	Delete(user entity.User) error
}

type Comments interface {
	Store(comment entity.Comment) error
	Fetch(postId int64) ([]entity.Comment, error)
	GetById(id int64) (entity.Comment, error)
	Update(comment entity.Comment) error
	Delete(post entity.Comment) error
	StoreLike(comment entity.Comment) error
	DeleteLike(comment entity.Comment) error
	StoreDislike(comment entity.Comment) error
	DeleteDislike(comment entity.Comment) error
	FetchReactions(id int64) (entity.Comment, error)
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
