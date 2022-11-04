package usecase

import (
	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

type PostsUseCase struct {
	repo sqlite.Posts

	userUseCase    sqlite.Users
	commentUseCase sqlite.Comments
}

func NewPostsUseCase(repo sqlite.Posts, userUseCase sqlite.Users, commentUseCase sqlite.Comments) *PostsUseCase {
	return &PostsUseCase{
		repo:           repo,
		userUseCase:    userUseCase,
		commentUseCase: commentUseCase,
	}
}

func (pu *PostsUseCase) CreatePost(p entity.Post) error {
	return nil
}

func (pu *PostsUseCase) GetAllPosts(p entity.Post) ([]entity.Post, error) {
	var posts []entity.Post
	return posts, nil
}

func (pu *PostsUseCase) GetOne(p entity.Post) (entity.Post, error) {
	var post entity.Post
	return post, nil
}

func (pu *PostsUseCase) UpdatePost(p entity.Post) (entity.Post, error) {
	var post entity.Post
	return post, nil
}

func (pu *PostsUseCase) DeletePost(p entity.Post) error {
	return nil
}

func (pu *PostsUseCase) MakeReaction(p entity.Post) error {
	return nil
}

func (pu *PostsUseCase) DeleteReaction(p entity.Post) error {
	return nil
}
