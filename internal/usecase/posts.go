package usecase

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"strings"
	"sync"
)

type PostsUseCase struct {
	repo repository.Posts

	userRepo    repository.Users
	commentRepo repository.Comments
}

const (
	PostCommentedQuery = "commented"
	PostLikedQuery     = "liked"
	PostDislikedQuery  = "disliked"
	ReactionLike       = "like"
	ReactionDislike    = "dislike"
	UniqueReactionErr  = "UNIQUE constraint failed"
	NoRowsResultErr    = "no rows in result set"
)

func NewPostsUseCase(repo repository.Posts, userUseCase repository.Users, commentUseCase repository.Comments) *PostsUseCase {
	return &PostsUseCase{
		repo:        repo,
		userRepo:    userUseCase,
		commentRepo: commentUseCase,
	}
}

func (pu *PostsUseCase) CreatePost(post entity.Post) error {
	post.Date = getRegTime(DateAndTimeFormat)
	err := pu.repo.Store(&post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreatePost #1 - %w", err)
	}
	err = pu.repo.StoreTopicReference(post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreatePost #2 - %w", err)
	}
	return nil
}

func (pu *PostsUseCase) GetAllPosts() ([]entity.Post, error) {
	posts, err := pu.repo.Fetch()
	if err != nil {
		return posts, fmt.Errorf("PostsUseCase - GetAllPosts #1 - %w", err)
	}
	if len(posts) == 0 {
		return posts, entity.ErrPostNotFound
	}
	err = pu.fillPostDetails(&posts)
	if err != nil {
		return posts, fmt.Errorf("PostsUseCase - GetAllPosts #2 - %w", err)
	}
	return posts, nil
}

func (pu *PostsUseCase) GetById(id int64) (entity.Post, error) {
	post, err := pu.repo.GetById(id)
	post.Id = id
	if err != nil {
		if strings.Contains(err.Error(), NoRowsResultErr) {
			return post, entity.ErrPostNotFound
		}
		return post, fmt.Errorf("PostsUseCase - GetById #1 - %w", err)
	}
	posts := []entity.Post{post}
	err = pu.fillPostDetails(&posts)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetById #2 - %w", err)
	}
	return posts[0], nil
}

func (pu *PostsUseCase) GetOneByCategory(category string) (entity.Post, error) {
	var post entity.Post
	ids, err := pu.repo.GetIdsByCategory(category)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetOneByCategory #1 - %w", err)
	}
	if len(ids) == 0 {
		return post, entity.ErrPostNotFound
	}
	post, err = pu.GetById(ids[0])
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetOneByCategory #2 - %w", err)
	}
	return post, nil
}

func (pu *PostsUseCase) GetAllByCategory(category string) ([]entity.Post, error) {
	var posts []entity.Post
	ids, err := pu.repo.GetIdsByCategory(category)
	if err != nil {
		return posts, fmt.Errorf("PostsUseCase - GetAllByCategory #1 - %w", err)
	}
	if len(ids) == 0 {
		return posts, entity.ErrPostNotFound
	}

	for i := 0; i < len(ids); i++ {
		post, err := pu.GetById(ids[i])
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetAllByCategory #2 - %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pu *PostsUseCase) GetAllCategories() ([]string, error) {
	categories, err := pu.repo.GetExistedCategories()
	if err != nil {
		return categories, fmt.Errorf("PostsUseCase - GetAllCategories #1 - %w", err)
	}

	if err != nil {
		return categories, fmt.Errorf("PostsUseCase - GetAllCategories #2 - %w", err)
	}
	return categories, nil
}

func (pu *PostsUseCase) UpdatePost(post entity.Post) error {
	err := pu.repo.Update(post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - UpdatePost #1 - %w", err)
	}
	return nil
}

func (pu *PostsUseCase) DeletePost(post entity.Post) error {
	err := pu.repo.Delete(post)
	if err != nil {
		return fmt.Errorf("PostsUseCase - UpdatePost #2 - %w", err)
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
					return fmt.Errorf("PostsUseCase - MakeReaction #1 - %w", err)
				}
				return nil
			}
			return fmt.Errorf("PostsUseCase - MakeReaction #2 - %w", err)
		}
		err = pu.repo.DeleteDislike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - MakeReaction #3 - %w", err)
		}
	case ReactionDislike:
		err := pu.repo.StoreDislike(post)
		if err != nil {
			if strings.Contains(err.Error(), UniqueReactionErr) {
				err = pu.repo.DeleteDislike(post)
				if err != nil {
					return fmt.Errorf("PostsUseCase - MakeReaction #4 - %w", err)
				}
				return nil
			}
			return fmt.Errorf("PostsUseCase - MakeReaction #5 - %w", err)
		}
		err = pu.repo.DeleteLike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - MakeReaction #6 - %w", err)
		}
	}
	return nil
}

func (pu *PostsUseCase) DeleteReaction(post entity.Post, command string) error {
	switch command {
	case ReactionLike:
		err := pu.repo.DeleteLike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - DeleteReaction #1 - %w", err)
		}
	case ReactionDislike:
		err := pu.repo.DeleteDislike(post)
		if err != nil {
			return fmt.Errorf("PostsUseCase - DeleteReaction #2 - %w", err)
		}
	}
	return nil
}

func (pu *PostsUseCase) fillPostDetails(posts *[]entity.Post) error {
	wgCategory := sync.WaitGroup{}
	wgComments := sync.WaitGroup{}

	categorySlice := make([][]string, len(*posts))
	commentsSlice := make([][]entity.Comment, len(*posts))

	errChan := make(chan error, 1)
	categoryDone := make(chan interface{})
	commentsDone := make(chan interface{})

	//collecting categories data
	for i := 0; i < len(*posts); i++ {
		wgCategory.Add(1)
		go func(n int) {
			defer wgCategory.Done()
			categories, err := pu.repo.GetRelatedCategories((*posts)[n])
			if err != nil {
				errChan <- fmt.Errorf("PostsUseCase - fillPostDetails #1 - %w", err)
			}
			categorySlice[n] = categories
		}(i)
	}

	//collecting comments data
	for i := 0; i < len(*posts); i++ {
		wgComments.Add(1)
		go func(n int) {
			defer wgComments.Done()
			comments, err := pu.commentRepo.Fetch((*posts)[n].Id)

			if err != nil {
				errChan <- fmt.Errorf("PostsUseCase - fillPostDetails #3 - %w", err)
			}
			for j := 0; j < len(comments); j++ {
				comments[j].User, err = pu.userRepo.GetById(comments[j].User.Id)
				if err != nil {
					continue
				}
			}
			commentsSlice[n] = comments
		}(i)
	}

	//filling categories data
	go func() {
		wgCategory.Wait()
		for i := 0; i < len(*posts); i++ {
			(*posts)[i].Categories = append((*posts)[i].Categories, categorySlice[i]...)
		}
		close(categoryDone)
	}()

	//filling comments data
	go func() {
		wgComments.Wait()
		for i := 0; i < len(*posts); i++ {
			(*posts)[i].Comments = append((*posts)[i].Comments, commentsSlice[i]...)
			(*posts)[i].TotalComments = int64(len((*posts)[i].Comments))
			if len(commentsSlice[i]) != 0 {
				(*posts)[i].LastComment = commentsSlice[i][len(commentsSlice[i])-1]
				(*posts)[i].LastCommentExist = true
			}

		}
		close(commentsDone)
	}()

	// if error occurs returning it
	for {
		select {
		case <-errChan:
			return <-errChan
		case _, ok := <-categoryDone:
			if !ok {
				categoryDone = nil
			}
		case _, ok := <-commentsDone:
			if !ok {
				commentsDone = nil
			}
		}
		if commentsDone == nil && categoryDone == nil {
			break
		}
	}

	return nil
}

func (pu *PostsUseCase) GetReactions(id int64) (entity.Post, error) {
	post, err := pu.repo.FetchReactions(id)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetReactions #1 - %w", err)
	}

	return post, nil
}

func (pu *PostsUseCase) CreateCategories(categories []string) error {
	err := pu.repo.StoreCategories(categories)
	if err != nil {
		return fmt.Errorf("PostsUseCase - GetReactions #2 - %w", err)
	}

	return nil
}
