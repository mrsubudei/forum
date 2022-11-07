package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"strings"
)

type PostsUseCase struct {
	repo repository.Posts

	userUseCase    repository.Users
	commentUseCase repository.Comments
}

const (
	PostCommentedQuery = "commented"
	PostLikedQuery     = "liked"
	PostDislikedQuery  = "disliked"
	ReactionLike       = "like"
	ReactionDislike    = "dislike"
	UniqueReactionErr  = "UNIQUE constraint failed"
)

func NewPostsUseCase(repo repository.Posts, userUseCase repository.Users, commentUseCase repository.Comments) *PostsUseCase {
	return &PostsUseCase{
		repo:           repo,
		userUseCase:    userUseCase,
		commentUseCase: commentUseCase,
	}
}

func (pu *PostsUseCase) CreatePost(post entity.Post) error {
	err := pu.repo.Store(&post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreatePost - %w", err)
	}
	err = pu.repo.StoreTopicReference(post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreatePost - %w", err)
	}
	return nil
}

func (pu *PostsUseCase) GetAllPosts(p entity.Post) ([]entity.Post, error) {
	var posts []entity.Post
	return posts, nil
}

func (pu *PostsUseCase) GetById(id int64) (entity.Post, error) {
	post, err := pu.repo.GetById(id)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetById - %w", err)
	}

	return post, nil
}

func (pu *PostsUseCase) GetByCategory(category string) (entity.Post, error) {
	var post entity.Post
	id, err := pu.repo.GetIdByCategory(category)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetByCategory #1 - %w", err)
	}
	post, err = pu.repo.GetById(id)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetByCategory #2 - %w", err)
	}
	return post, nil
}

func (pu *PostsUseCase) UpdatePost(post entity.Post) (entity.Post, error) {
	err := pu.repo.Update(post)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - UpdatePost - %w", err)
	}
	return post, nil
}

func (pu *PostsUseCase) DeletePost(post entity.Post) error {
	err := pu.repo.Delete(post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - DeletePost - %w", err)
	}
	return nil
}

func (pu *PostsUseCase) MakeReaction(post entity.Post, command string) error {
	switch command {
	case ReactionLike:
		err := pu.repo.StoreLike(post)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = pu.repo.DeleteLike(post)
				if err != nil {
					return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionLike - %w", err)
				}
				return nil
			}
			return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionLike -  %w", err)
		}
		err = pu.repo.DeleteDislike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionLike -  %w", err)
		}
	case ReactionDislike:
		err := pu.repo.StoreDislike(post)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = pu.repo.DeleteDislike(post)
				if err != nil {
					return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionDisike - %w", err)
				}
				return nil
			}
			return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionDisike - %w", err)
		}
		err = pu.repo.DeleteLike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - MakeReaction - case ReactionDisike - %w", err)
		}
	}
	return nil
}

func (pu *PostsUseCase) DeleteReaction(post entity.Post, command string) error {
	switch command {
	case ReactionLike:
		err := pu.repo.DeleteLike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - DeleteReaction - %w", err)
		}
	case ReactionDislike:
		err := pu.repo.DeleteDislike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - DeleteReaction - %w", err)
		}
	}
	return nil
}

func (pu *PostsUseCase) fillPostDetails(posts *[]entity.Post) error {
	// for i := 0; i < len(*posts); i++ {
	// 	categories, err := pu.repo.GetRelatedCategories((*posts)[i])
	// 	if err != nil {
	// 		return fmt.Errorf("PostsUseCase - fillPostDetails - %w", err)
	// 	}
	// 	(*posts)[i].Category = categories
	// }
	// mapPosts := make(map[int64]entity.Post)
	// for _, post := range *posts {
	// 	mapPosts[post.Id] = entity.Post{}
	// }
	// chanPost := make(chan entity.Post)
	// for postId :=
	return nil
}
