package usecase

import (
	"forum/internal/entity"
	"forum/internal/repository"
	"forum/pkg/auth"
	"forum/pkg/hasher"
)

type Posts interface {
	CreatePost(p entity.Post) error
	GetAllPosts(p entity.Post) ([]entity.Post, error)
	GetById(id int64) (entity.Post, error)
	GetByCategory(category string) (entity.Post, error)
	UpdatePost(p entity.Post) (entity.Post, error)
	DeletePost(p entity.Post) error
	MakeReaction(p entity.Post, command string) error
	DeleteReaction(post entity.Post, command string) error
}

type Users interface {
	SignUp(u entity.User) error
	SignIn(u entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetById(id int64) (entity.User, error)
	GetSession(id int64) (entity.User, error)
	CheckSession(u entity.User) (bool, error)
	UpdateUserInfo(u entity.User, query string) error
	UpdateSession(u entity.User) error
	DeleteUser(u entity.User) error
}

type Comments interface {
	WriteComment(c entity.Comment) error
	GetAllComments() ([]entity.Comment, error)
	UpdateComment(c entity.Comment) error
	DeleteComment(c entity.Comment) error
	MakeReaction(comment entity.Comment, command string) error
	DeleteReaction(c entity.Comment) error
}

type UseCases struct {
	Posts    Posts
	Users    Users
	Comments Comments
}

type Dependencies struct {
	Repos        *repository.Repositories
	Hasher       hasher.PasswordHasher
	TokenManager auth.TokenManager
}

func NewUseCases(deps Dependencies) *UseCases {
	postUseCase := NewPostsUseCase(deps.Repos.Posts, deps.Repos.Users, deps.Repos.Comments)
	userUseCase := NewUsersUseCase(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.Repos.Posts, deps.Repos.Comments)
	commentsUseCase := NewCommentUseCase(deps.Repos.Comments, deps.Repos.Posts, deps.Repos.Users)

	return &UseCases{
		Posts:    postUseCase,
		Users:    userUseCase,
		Comments: commentsUseCase,
	}
}
