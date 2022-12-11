package usecase

import (
	"fmt"
	"strings"
	"sync"

	"forum/internal/entity"
	"forum/internal/repository"
)

type PostsUseCase struct {
	repo repository.Posts

	userRepo    repository.Users
	commentRepo repository.Comments
}

const (
	PostCommentedQuery = "commented"
	PostAuthorQuery    = "author"
	PostLikedQuery     = "liked"
	PostDislikedQuery  = "disliked"
	ReactionLike       = "like"
	ReactionDislike    = "dislike"
	UniqueReactionErr  = "UNIQUE constraint failed"
	NoRowsResultErr    = "no rows in result set"
)

func NewPostsUseCase(repo repository.Posts, usersRepo repository.Users, commentsRepo repository.Comments) *PostsUseCase {
	return &PostsUseCase{
		repo:        repo,
		userRepo:    usersRepo,
		commentRepo: commentsRepo,
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

	if len(posts) != 0 {
		err = pu.fillPostDetails(&posts)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetAllPosts #2 - %w", err)
		}
	}
	return posts, nil
}

func (pu *PostsUseCase) GetPostsByQuery(user entity.User, query string) ([]entity.Post, error) {
	var posts []entity.Post
	var err error
	switch query {
	case PostAuthorQuery:
		posts, err = pu.repo.FetchByAuthor(user)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #1 - %w", err)
		}
	case PostLikedQuery:
		ids, err := pu.repo.FetchIdsByReaction(user, PostLikedQuery)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #2 - %w", err)
		}
		for i := 0; i < len(ids); i++ {
			post, err := pu.repo.GetById(ids[i])
			if err != nil {
				return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #3 - %w", err)
			}
			posts = append(posts, post)
		}
	case PostDislikedQuery:
		ids, err := pu.repo.FetchIdsByReaction(user, PostDislikedQuery)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #4 - %w", err)
		}
		for i := 0; i < len(ids); i++ {
			post, err := pu.repo.GetById(ids[i])
			if err != nil {
				return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #5 - %w", err)
			}
			posts = append(posts, post)
		}
	case PostCommentedQuery:
		ids, err := pu.commentRepo.GetPostIds(user)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #5 - %w", err)
		}
		for i := 0; i < len(ids); i++ {
			post, err := pu.repo.GetById(ids[i])
			if err != nil {
				return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #6 - %w", err)
			}
			posts = append(posts, post)
		}
	}

	if len(posts) != 0 {
		err = pu.fillPostDetails(&posts)
		if err != nil {
			return posts, fmt.Errorf("PostsUseCase - GetPostsByQuery #3 - %w", err)
		}
	}
	return posts, nil
}

func (pu *PostsUseCase) GetById(id int64) (entity.Post, error) {
	post, err := pu.repo.GetById(id)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetById #1 - %w", err)
	}
	user, err := pu.userRepo.GetById(post.User.Id)
	if user.Gender == UserGenderMale {
		user.Male = true
	} else if user.Gender == UserGenderFemale {
		user.Female = true
	}
	post.User = user

	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetById #2 - %w", err)
	}
	post.Id = id
	if err != nil {
		if strings.Contains(err.Error(), NoRowsResultErr) {
			return post, entity.ErrPostNotFound
		}
		return post, fmt.Errorf("PostsUseCase - GetById #3 - %w", err)
	}
	posts := []entity.Post{post}
	err = pu.fillPostDetails(&posts)
	if err != nil {
		return post, fmt.Errorf("PostsUseCase - GetById #4 - %w", err)
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
		fmt.Println(err)
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
		return categories, fmt.Errorf("PostsUseCase - GetAllCategories - %w", err)
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

	// collecting categories data
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

	// collecting comments data
	for i := 0; i < len(*posts); i++ {
		wgComments.Add(1)
		go func(n int) {
			defer wgComments.Done()
			comments, err := pu.commentRepo.Fetch((*posts)[n].Id)
			if err != nil {
				errChan <- fmt.Errorf("PostsUseCase - fillPostDetails #3 - %w", err)
			}
			for j := 0; j < len(comments); j++ {
				comments[j].ContentWeb = strings.Split(comments[j].Content, "\\n")
				comments[j].User, err = pu.userRepo.GetById(comments[j].User.Id)
				if comments[j].User.Gender == UserGenderMale {
					comments[j].User.Male = true
				} else if comments[j].User.Gender == UserGenderFemale {
					comments[j].User.Female = true
				}
				if err != nil {
					continue
				}
			}
			commentsSlice[n] = comments
		}(i)
	}

	// filling categories data
	go func() {
		wgCategory.Wait()
		for i := 0; i < len(*posts); i++ {
			(*posts)[i].Categories = append((*posts)[i].Categories, categorySlice[i]...)
		}
		close(categoryDone)
	}()

	// filling comments data
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

func (pu *PostsUseCase) GetReactions(id int64, query string) ([]entity.User, error) {
	var users []entity.User
	post, err := pu.repo.FetchReactions(id)
	if err != nil {
		return users, fmt.Errorf("PostsUseCase - GetReactions #1 - %w", err)
	}
	switch query {
	case PostLikedQuery:
		for i := 0; i < len(post.Likes); i++ {
			user, err := pu.userRepo.GetById(post.Likes[i].UserId)
			if err != nil {
				return users, fmt.Errorf("PostsUseCase - GetReactions #2 - %w", err)
			}
			users = append(users, user)
		}
	case PostDislikedQuery:
		for i := 0; i < len(post.Dislikes); i++ {
			user, err := pu.userRepo.GetById(post.Dislikes[i].UserId)
			if err != nil {
				return users, fmt.Errorf("PostsUseCase - GetReactions #3 - %w", err)
			}
			users = append(users, user)
		}
	}

	return users, nil
}

func (pu *PostsUseCase) CreateCategories(categories []string) error {
	existed, err := pu.repo.GetExistedCategories()
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreateCategories #1 - %w", err)
	}

	var categoriesToAdd []string
	for i := 0; i < len(categories); i++ {
		exist := false
		for j := 0; j < len(existed); j++ {
			if categories[i] == existed[j] {
				exist = true
			}
		}
		if !exist {
			categoriesToAdd = append(categoriesToAdd, categories[i])
		}
	}

	err = pu.repo.StoreCategories(categoriesToAdd)
	if err != nil {
		return fmt.Errorf("PostsUseCase - CreateCategories #2 - %w", err)
	}

	return nil
}
