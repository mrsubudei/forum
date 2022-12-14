package mock_repository

import (
	"fmt"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
)

var errNoRows = fmt.Errorf("no rows in result set")

type MockRepos struct {
	Users    *UsersMockRepo
	Posts    *PostsMockRepo
	Comments *CommentsMockRepo
}

func NewMockRepos() *MockRepos {
	return &MockRepos{
		Users:    NewUsersMockRepo(),
		Posts:    NewPostsMockrepo(),
		Comments: NewCommentsMockrepo(),
	}
}

type UsersMockRepo struct {
	users []entity.User
}

func NewUsersMockRepo() *UsersMockRepo {
	return &UsersMockRepo{}
}

func (um *UsersMockRepo) Store(user entity.User) error {
	for _, v := range um.users {
		if v.Name == user.Name {
			return fmt.Errorf(usecase.UniqueNameErr)
		} else if v.Email == user.Email {
			return fmt.Errorf(usecase.UniqueEmailErr)
		}
	}
	um.users = append(um.users, user)
	return nil
}

func (um *UsersMockRepo) Fetch() ([]entity.User, error) {
	return um.users, nil
}

func (um *UsersMockRepo) GetId(user entity.User) (int64, error) {
	var id int64
	switch {
	case user.SessionToken != "":
		for _, v := range um.users {
			if v.SessionToken == user.SessionToken {
				return v.Id, nil
			}
		}
	case user.Name != "":
		for _, v := range um.users {
			if v.Name == user.Name {
				return v.Id, nil
			}
		}
	case user.Email != "":
		for _, v := range um.users {
			if v.Email == user.Email {
				return v.Id, nil
			}
		}
	}

	return id, errNoRows
}

func (um *UsersMockRepo) GetById(n int64) (entity.User, error) {
	for _, v := range um.users {
		if v.Id == n {
			return v, nil
		}
	}
	return entity.User{}, errNoRows
}

func (um *UsersMockRepo) GetSession(n int64) (entity.User, error) {
	for _, v := range um.users {
		if v.Id == n {
			return v, nil
		}
	}
	return entity.User{}, errNoRows
}

func (um *UsersMockRepo) UpdateInfo(user entity.User) error {
	for i := 0; i < len(um.users); i++ {
		if um.users[i].Id == user.Id {
			um.users[i] = user
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) UpdatePassword(user entity.User) error {
	for i := 0; i < len(um.users); i++ {
		if um.users[i].Id == user.Id {
			um.users[i].Password = user.Password
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) NewSession(user entity.User) error {
	for i := 0; i < len(um.users); i++ {
		if um.users[i].Id == user.Id {
			um.users[i].SessionToken = user.SessionToken
			um.users[i].SessionTTL = user.SessionTTL
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) UpdateSession(user entity.User) error {
	for i := 0; i < len(um.users); i++ {
		if um.users[i].Id == user.Id {
			um.users[i].SessionTTL = user.SessionTTL
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) Delete(user entity.User) error {
	newUsers := []entity.User{}
	found := false
	for _, v := range um.users {
		if v.Id != user.Id {
			newUsers = append(newUsers, v)
		} else {
			found = true
		}
	}
	um.users = newUsers

	if found {
		return nil
	} else {
		return errNoRows
	}
}

type PostsMockRepo struct {
	posts     []entity.Post
	allTopics map[string]bool
}

func NewPostsMockrepo() *PostsMockRepo {
	return &PostsMockRepo{}
}

func (pm *PostsMockRepo) Store(post *entity.Post) error {
	pm.posts = append(pm.posts, *post)
	return nil
}

func (pm *PostsMockRepo) Fetch() ([]entity.Post, error) {
	return pm.posts, nil
}

func (pm *PostsMockRepo) FetchByAuthor(user entity.User) ([]entity.Post, error) {
	posts := []entity.Post{}
	for i := 0; i < len(pm.posts); i++ {
		if pm.posts[i].User.Id == user.Id {
			posts = append(posts, pm.posts[i])
		}
	}
	return posts, nil
}

func (pm *PostsMockRepo) GetById(id int64) (entity.Post, error) {
	for i := 0; i < len(pm.posts); i++ {
		if pm.posts[i].Id == id {
			return pm.posts[i], nil
		}
	}
	return entity.Post{}, nil
}

func (pm *PostsMockRepo) GetIdsByCategory(category string) ([]int64, error) {
	var ids []int64
	for _, val := range pm.posts {
		for _, v := range val.Categories {
			if v == category {
				ids = append(ids, val.Id)
			}
		}
	}
	return ids, nil
}

func (pm *PostsMockRepo) FetchIdsByReaction(user entity.User, reaction string) ([]int64, error) {
	var ids []int64
	switch reaction {
	case sqlite.QueryLiked:
		for _, val := range pm.posts {
			for _, v := range val.Likes {
				if v.UserId == user.Id {
					ids = append(ids, val.Id)
				}
			}
		}
	case sqlite.QueryDislike:
		for _, val := range pm.posts {
			for _, v := range val.Dislikes {
				if v.UserId == user.Id {
					ids = append(ids, val.Id)
				}
			}
		}
	}
	return ids, nil
}

func (pm *PostsMockRepo) Update(post entity.Post) error {
	for i := 0; i < len(pm.posts); i++ {
		if pm.posts[i].Id == post.Id {
			pm.posts[i] = post
			return nil
		}
	}
	return errNoRows
}

func (pm *PostsMockRepo) Delete(post entity.Post) error {
	newPosts := []entity.Post{}
	found := false
	for _, v := range pm.posts {
		if v.Id != post.Id {
			newPosts = append(newPosts, v)
		} else {
			found = true
		}
	}
	pm.posts = newPosts

	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) StoreLike(post entity.Post) error {
	found := false
	like := entity.Reaction{UserId: post.User.Id}
	for i := 0; i < len(pm.posts); i++ {
		if pm.posts[i].Id == post.Id {
			pm.posts[i].Likes = append(pm.posts[i].Likes, like)
			found = true
		}
	}
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) StoreDislike(post entity.Post) error {
	found := false
	dislike := entity.Reaction{UserId: post.User.Id}
	for i := 0; i < len(pm.posts); i++ {
		if pm.posts[i].Id == post.Id {
			pm.posts[i].Dislikes = append(pm.posts[i].Dislikes, dislike)
			found = true
		}
	}
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) DeleteLike(post entity.Post) error {
	found := false
	newPosts := []entity.Post{}
	for _, val := range pm.posts {
		for _, v := range val.Likes {
			if v.UserId != post.User.Id {
				newPosts = append(newPosts, val)
			} else {
				found = true
			}
		}
	}
	pm.posts = newPosts
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) DeleteDislike(post entity.Post) error {
	found := false
	newPosts := []entity.Post{}
	for _, val := range pm.posts {
		for _, v := range val.Dislikes {
			if v.UserId != post.User.Id {
				newPosts = append(newPosts, val)
			} else {
				found = true
			}
		}
	}
	pm.posts = newPosts
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) StoreTopicReference(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) GetRelatedCategories(post entity.Post) ([]string, error) {
	for _, v := range pm.posts {
		if v.Id == post.Id {
			return v.Categories, nil
		}
	}
	return nil, errNoRows
}

func (pm *PostsMockRepo) FetchReactions(id int64) (entity.Post, error) {
	for _, v := range pm.posts {
		if v.Id == id {
			return v, nil
		}
	}
	return entity.Post{}, errNoRows
}

func (pm *PostsMockRepo) StoreCategories(categories []string) error {
	for _, v := range categories {
		pm.allTopics[v] = true
	}
	return nil
}

func (pm *PostsMockRepo) GetExistedCategories() ([]string, error) {
	categories := []string{}
	for key := range pm.allTopics {
		categories = append(categories, key)
	}
	return categories, nil
}

type CommentsMockRepo struct {
	posts []entity.Comment
}

func NewCommentsMockrepo() *CommentsMockRepo {
	return &CommentsMockRepo{}
}

func (cm *CommentsMockRepo) Store(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) Fetch(postId int64) ([]entity.Comment, error) {
	comments := []entity.Comment{}
	return comments, nil
}

func (cm *CommentsMockRepo) GetById(id int64) (entity.Comment, error) {
	comment := entity.Comment{}
	return comment, nil
}

func (cm *CommentsMockRepo) GetPostIds(user entity.User) ([]int64, error) {
	var ids []int64
	return ids, nil
}

func (cm *CommentsMockRepo) Update(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) Delete(post entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) StoreLike(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) DeleteLike(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) StoreDislike(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) DeleteDislike(comment entity.Comment) error {
	return nil
}

func (cm *CommentsMockRepo) FetchReactions(id int64) (entity.Comment, error) {
	comment := entity.Comment{}
	return comment, nil
}
