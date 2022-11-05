package usecase

import (
	"forum/internal/entity"
	"forum/internal/repository"
)

type CommentsUseCase struct {
	repo repository.Comments

	postUseCase repository.Posts
	userUseCase repository.Users
}

func NewCommentUseCase(repo repository.Comments, postUseCase repository.Posts, userUseCase repository.Users) *CommentsUseCase {
	return &CommentsUseCase{
		repo:        repo,
		postUseCase: postUseCase,
		userUseCase: userUseCase,
	}
}

func (cu *CommentsUseCase) WriteComment(p entity.Post) error {
	return nil
}

func (cu *CommentsUseCase) GetAllComments(p entity.Post) error {
	return nil
}

func (cu *CommentsUseCase) UpdateComment(c entity.Comment) error {
	return nil
}

func (cu *CommentsUseCase) DeleteComment(c entity.Comment) error {
	return nil
}

func (cu *CommentsUseCase) MakeReaction(c entity.Comment) error {
	return nil
}

func (cu *CommentsUseCase) DeleteReaction(c entity.Comment) error {
	return nil
}
