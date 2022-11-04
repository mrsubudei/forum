package usecase

import (
	"forum/internal/entity"
	"forum/internal/repository/sqlite"
)

type CommentsUseCase struct {
	repo sqlite.Comments

	postUseCase sqlite.Posts
	userUseCase sqlite.Users
}

func NewCommentUseCase(repo sqlite.Comments, postUseCase sqlite.Posts, userUseCase sqlite.Users) *CommentsUseCase {
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
