package usecase

import "forum/internal/entity"

type (
	Communication interface {
		GetAllUsers() ([]entity.User, error)
		GetUser(userId int) (entity.User, error)
		GetUserPostIds(userId int) ([]int, error)
		GetAllPosts() ([]entity.Post, error)
		CreateDB() error
		CreateUser(u *entity.User) (int, error)
		CreatePost(p *entity.Post) (int, error)
		CreateComment(p *entity.Post, date, content string) error
		PutPostLike(u *entity.User, postId int, date string) error
		PutPostDisLike(u *entity.User, postId int, date string) error
		PutCommentLike(u *entity.User, commentId int, date string) error
		PutCommentDisLike(u *entity.User, commentId int, date string) error
		CreateTopics(name []string) error
		CreatePostRef(p *entity.Post, name []string) error
	}

	CommunicationRepo interface {
		GetAllUsers() ([]entity.User, error)
		GetUser(userId int) (entity.User, error)
		GetUserPostIds(userId int) ([]int, error)
		GetAllPosts() ([]entity.Post, error)
		CreateDB() error
		CreateUser(u *entity.User) (int, error)
		CreatePost(p *entity.Post) (int, error)
		CreateComment(p *entity.Post, date, content string) error
		PutPostLike(u *entity.User, postId int, date string) error
		PutPostDisLike(u *entity.User, postId int, date string) error
		PutCommentLike(u *entity.User, commentId int, date string) error
		PutCommentDisLike(u *entity.User, commentId int, date string) error
		CreateTopics(name []string) error
		CreatePostRef(p *entity.Post, name []string) error
	}
)
