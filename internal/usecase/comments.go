package usecase

import (
	"fmt"
	"strings"

	"forum/internal/entity"
	"forum/internal/repository"
)

type CommentsUseCase struct {
	repo repository.Comments

	postRepo repository.Posts
	userRepo repository.Users
}

func NewCommentsUseCase(repo repository.Comments, postsRepo repository.Posts, usersRepo repository.Users) *CommentsUseCase {
	return &CommentsUseCase{
		repo:     repo,
		postRepo: postsRepo,
		userRepo: usersRepo,
	}
}

func (cu *CommentsUseCase) WriteComment(comment entity.Comment) error {
	comment.Date = getRegTime(DateAndTimeFormat)
	err := cu.repo.Store(comment)
	if err != nil {
		return fmt.Errorf("CommentsUseCase - WriteComment - %w", err)
	}

	return nil
}

func (cu *CommentsUseCase) GetAllComments(postId int64) ([]entity.Comment, error) {
	comments, err := cu.repo.Fetch(postId)
	if err != nil {
		return nil, fmt.Errorf("CommentsUseCase - GetAllComments #1 - %w", err)
	}

	for i := 0; i < len(comments); i++ {
		user, err := cu.userRepo.GetById(comments[i].User.Id)
		if err != nil {
			return nil, fmt.Errorf("CommentsUseCase - GetAllComments #2 - %w", err)
		}
		comments[i].User = user
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
					return fmt.Errorf("CommentsUseCase - MakeReaction #1 - %w", err)
				}
				return nil
			}
			return fmt.Errorf("CommentsUseCase - MakeReaction #2 -  %w", err)
		}
		err = cu.repo.DeleteDislike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - MakeReaction #3 -  %w", err)
		}
	case ReactionDislike:
		err := cu.repo.StoreDislike(comment)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = cu.repo.DeleteDislike(comment)
				if err != nil {
					return fmt.Errorf("CommentsUseCase - MakeReaction #4 - %w", err)
				}
				return nil
			}
			return fmt.Errorf("CommentsUseCase - MakeReaction #5 - %w", err)
		}
		err = cu.repo.DeleteLike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - MakeReaction #6 - %w", err)
		}
	}
	return nil
}

func (cu *CommentsUseCase) DeleteReaction(comment entity.Comment, command string) error {
	switch command {
	case ReactionLike:
		err := cu.repo.DeleteLike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - DeleteReaction #1 - %w", err)
		}
	case ReactionDislike:
		err := cu.repo.DeleteDislike(comment)
		if err != nil {
			return fmt.Errorf("CommentsUseCase - DeleteReaction #2 - %w", err)
		}
	}
	return nil
}

func (cu *CommentsUseCase) GetReactions(id int64, query string) ([]entity.User, error) {
	var users []entity.User
	comment, err := cu.repo.FetchReactions(id)
	if err != nil {
		return users, fmt.Errorf("PostsUseCase - GetReactions #1 - %w", err)
	}

	switch query {
	case PostLikedQuery:
		for i := 0; i < len(comment.Likes); i++ {
			user, err := cu.userRepo.GetById(comment.Likes[i].UserId)
			if err != nil {
				return users, fmt.Errorf("PostsUseCase - GetReactions #2 - %w", err)
			}
			users = append(users, user)
		}
	case PostDislikedQuery:
		for i := 0; i < len(comment.Dislikes); i++ {
			user, err := cu.userRepo.GetById(comment.Dislikes[i].UserId)
			if err != nil {
				return users, fmt.Errorf("PostsUseCase - GetReactions #3 - %w", err)
			}
			users = append(users, user)
		}
	}

	return users, nil
}
