package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
)

type PostsUseCase struct {
	repo repository.Posts

	userUseCase    repository.Users
	commentUseCase repository.Comments
}

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
