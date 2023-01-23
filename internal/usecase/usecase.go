package usecase

import (
	"forum/internal/entity"
)

type Posts interface {
	CreatePost(p entity.Post) error
	GetAllPosts() ([]entity.Post, error)
	GetPostsByQuery(user entity.User, query string) ([]entity.Post, error)
	GetById(id int64) (entity.Post, error)
	GetAllByCategory(category string) ([]entity.Post, error)
	UpdatePost(post entity.Post) error
	DeletePost(p entity.Post) error
	MakeReaction(p entity.Post, command string) error
	DeleteReaction(post entity.Post, command string) error
	CreateCategories(categories []string) error
	GetAllCategories() ([]string, error)
	GetReactions(id int64, query string) ([]entity.User, error)
}

type Users interface {
	SignUp(u entity.User) error
	SignIn(u entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetById(id int64) (entity.User, error)
	GetIdBy(user entity.User) (int64, error)
	GetSession(id int64) (entity.User, error)
	CheckSession(u entity.User) (bool, error)
	UpdateUserInfo(u entity.User, query string) error
	UpdateSession(u entity.User) error
	DeleteSession(user entity.User) error
	DeleteUser(u entity.User) error
}

type Comments interface {
	WriteComment(c entity.Comment) error
	GetAllComments(postId int64) ([]entity.Comment, error)
	UpdateComment(c entity.Comment) error
	DeleteComment(c entity.Comment) error
	MakeReaction(c entity.Comment, command string) error
	DeleteReaction(c entity.Comment, command string) error
	GetReactions(id int64, query string) ([]entity.User, error)
}

type UseCases struct {
	Posts    Posts
	Users    Users
	Comments Comments
}

func NewUseCases(posts Posts, users Users, comments Comments) *UseCases {
	return &UseCases{
		Posts:    posts,
		Users:    users,
		Comments: comments,
	}
}
