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
	GetOne(p entity.Post) (entity.Post, error)
	UpdatePost(p entity.Post) (entity.Post, error)
	DeletePost(p entity.Post) error
	MakeReaction(p entity.Post) error
	DeleteReaction(p entity.Post) error
}

type Users interface {
	SignUp(u entity.User) error
	SignIn(u entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetById(id int) (entity.User, error)
	GetSession(id int) (entity.User, error)
	UpdateUserInfo(user entity.User, query string) error
	UpdateSession(user entity.User) error
	DeleteUser(u entity.User) error
}

type Comments interface {
	WriteComment(p entity.Post) error
	GetAllComments(p entity.Post) error
	UpdateComment(c entity.Comment) error
	DeleteComment(c entity.Comment) error
	MakeReaction(c entity.Comment) error
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
