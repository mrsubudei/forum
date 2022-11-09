package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"strings"
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

func (cu *CommentsUseCase) WriteComment(comment entity.Comment) error {
	err := cu.repo.Store(comment)
	if err != nil {
		return fmt.Errorf("CommentsUseCase - WriteComment - %w", err)
	}
	return nil
}

func (cu *CommentsUseCase) GetAllComments(postId int64) ([]entity.Comment, error) {
	comments, err := cu.repo.Fetch(postId)
	if err != nil {
		return nil, fmt.Errorf("CommentsUseCase - GetAllComments - %w", err)
	}
	return comments, nil
}

func (cu *CommentsUseCase) UpdateComment(comment entity.Comment) error {
	err := cu.repo.Update(comment)
	if err != nil {
		return fmt.Errorf("CommentsUseCase - UpdateComment - %w", err)
	}
	return nil
}

func (cu *CommentsUseCase) DeleteComment(comment entity.Comment) error {
	err := cu.repo.Delete(comment)
	if err != nil {
		return fmt.Errorf("CommentsUseCase - DeleteComment - %w", err)
	}
	return nil
}

func (cu *CommentsUseCase) MakeReaction(comment entity.Comment, command string) error {
	switch command {
	case ReactionLike:
		err := cu.repo.StoreLike(comment)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = cu.repo.DeleteLike(comment)
				if err != nil {
					return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionLike - %w", err)
				}
				return nil
			}
			return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionLike -  %w", err)
		}
		err = cu.repo.DeleteDislike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionLike -  %w", err)
		}
	case ReactionDislike:
		err := cu.repo.StoreDislike(comment)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = cu.repo.DeleteDislike(comment)
				if err != nil {
					return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionDisike - %w", err)
				}
				return nil
			}
			return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionDisike - %w", err)
		}
		err = cu.repo.DeleteLike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - MakeReaction - case ReactionDisike - %w", err)
		}
	}
	return nil
}

func (cu *CommentsUseCase) DeleteReaction(comment entity.Comment, command string) error {
	switch command {
	case ReactionLike:
		err := cu.repo.DeleteLike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - DeleteReaction - %w", err)
		}
	case ReactionDislike:
		err := cu.repo.DeleteDislike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - DeleteReaction - %w", err)
		}
	}
	return nil
}
