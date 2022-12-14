package mock_repository

import (
	"fmt"
	"time"

	"forum/internal/entity"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
)

var (
	errNoRows           = fmt.Errorf("no rows in result set")
	errUniqueConstraint = fmt.Errorf("UNIQUE constraint failed")
)

const DateAndTimeFormat = "2006-01-02 15:04:05"

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
	Users []entity.User
}

func NewUsersMockRepo() *UsersMockRepo {
	return &UsersMockRepo{}
}

func (um *UsersMockRepo) Store(user entity.User) error {
	for _, v := range um.Users {
		if v.Name == user.Name {
			return fmt.Errorf(usecase.UniqueNameErr)
		} else if v.Email == user.Email {
			return fmt.Errorf(usecase.UniqueEmailErr)
		}
	}
	um.Users = append(um.Users, user)
	return nil
}

func (um *UsersMockRepo) Fetch() ([]entity.User, error) {
	return um.Users, nil
}

func (um *UsersMockRepo) GetId(user entity.User) (int64, error) {
	var id int64
	switch {
	case user.SessionToken != "":
		for _, v := range um.Users {
			if v.SessionToken == user.SessionToken {
				return v.Id, nil
			}
		}
	case user.Name != "":
		for _, v := range um.Users {
			if v.Name == user.Name {
				return v.Id, nil
			}
		}
	case user.Email != "":
		for _, v := range um.Users {
			if v.Email == user.Email {
				return v.Id, nil
			}
		}
	}

	return id, errNoRows
}

func (um *UsersMockRepo) GetById(n int64) (entity.User, error) {
	for _, v := range um.Users {
		if v.Id == n {
			return v, nil
		}
	}
	return entity.User{}, errNoRows
}

func (um *UsersMockRepo) GetSession(n int64) (entity.User, error) {
	for _, v := range um.Users {
		if v.Id == n {
			return v, nil
		}
	}
	return entity.User{}, errNoRows
}

func (um *UsersMockRepo) UpdateInfo(user entity.User) error {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Id == user.Id {
			um.Users[i] = user
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) UpdatePassword(user entity.User) error {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Id == user.Id {
			um.Users[i].Password = user.Password
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) NewSession(user entity.User) error {
	ttl := user.SessionTTL.Format(DateAndTimeFormat)
	formatted, err := time.Parse(DateAndTimeFormat, ttl)
	if err != nil {
		return err
	}
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Id == user.Id {
			um.Users[i].SessionToken = user.SessionToken
			um.Users[i].SessionTTL = formatted
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) UpdateSession(user entity.User) error {
	for i := 0; i < len(um.Users); i++ {
		if um.Users[i].Id == user.Id {
			um.Users[i].SessionTTL = user.SessionTTL
			return nil
		}
	}
	return errNoRows
}

func (um *UsersMockRepo) Delete(user entity.User) error {
	newUsers := []entity.User{}
	found := false
	for i, v := range um.Users {
		if v.Id == user.Id {
			found = true
			newUsers = deleteElement(um.Users, i)
		}
	}
	if found {
		um.Users = []entity.User{}
		um.Users = append(um.Users, newUsers...)
		return nil
	} else {
		return errNoRows
	}
}

type PostsMockRepo struct {
	Posts     []entity.Post
	AllTopics map[string]bool
}

func NewPostsMockrepo() *PostsMockRepo {
	return &PostsMockRepo{}
}

func (pm *PostsMockRepo) Store(post *entity.Post) error {
	pm.Posts = append(pm.Posts, *post)
	return nil
}

func (pm *PostsMockRepo) Fetch() ([]entity.Post, error) {
	return pm.Posts, nil
}

func (pm *PostsMockRepo) FetchByAuthor(user entity.User) ([]entity.Post, error) {
	posts := []entity.Post{}
	for i := 0; i < len(pm.Posts); i++ {
		if pm.Posts[i].User.Id == user.Id {
			posts = append(posts, pm.Posts[i])
		}
	}
	return posts, nil
}

func (pm *PostsMockRepo) GetById(id int64) (entity.Post, error) {
	for i := 0; i < len(pm.Posts); i++ {
		if pm.Posts[i].Id == id {
			return pm.Posts[i], nil
		}
	}
	return entity.Post{}, errNoRows
}

func (pm *PostsMockRepo) GetIdsByCategory(category string) ([]int64, error) {
	var ids []int64
	for _, val := range pm.Posts {
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
		for _, val := range pm.Posts {
			for _, v := range val.Likes {
				if v.UserId == user.Id {
					ids = append(ids, val.Id)
				}
			}
		}
	case sqlite.QueryDislike:
		for _, val := range pm.Posts {
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
	for i := 0; i < len(pm.Posts); i++ {
		if pm.Posts[i].Id == post.Id {
			pm.Posts[i] = post
			return nil
		}
	}
	return errNoRows
}

func (pm *PostsMockRepo) Delete(post entity.Post) error {
	newPosts := []entity.Post{}
	found := false
	for i, v := range pm.Posts {
		if v.Id == post.Id {
			found = true
			newPosts = deleteElement(pm.Posts, i)
		}
	}
	if found {
		pm.Posts = []entity.Post{}
		pm.Posts = append(pm.Posts, newPosts...)
		return nil
	} else {
		return errNoRows
	}
}

func (pm *PostsMockRepo) StoreLike(post entity.Post) error {
	found := false
	like := entity.Reaction{UserId: post.User.Id}
	for i := 0; i < len(pm.Posts); i++ {
		if pm.Posts[i].Id == post.Id {
			alreadyReacted := false
			for _, v := range pm.Posts[i].Likes {
				if v.UserId == post.User.Id {
					alreadyReacted = true
				}
			}
			if alreadyReacted {
				return errUniqueConstraint
			}
			pm.Posts[i].Likes = append(pm.Posts[i].Likes, like)
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
	for i := 0; i < len(pm.Posts); i++ {
		if pm.Posts[i].Id == post.Id {
			alreadyReacted := false
			for _, v := range pm.Posts[i].Dislikes {
				if v.UserId == post.User.Id {
					alreadyReacted = true
				}
			}
			if alreadyReacted {
				return errUniqueConstraint
			}
			pm.Posts[i].Dislikes = append(pm.Posts[i].Dislikes, dislike)
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
	newLikes := []entity.Reaction{}
	idx := 0
	found := false
	for y, val := range pm.Posts {
		for i, v := range val.Likes {
			if v.UserId == post.User.Id {
				found = true
				newLikes = deleteElement(val.Likes, i)
				idx = y
				break
			}
		}
	}

	if found {
		pm.Posts[idx].Likes = []entity.Reaction{}
		pm.Posts[idx].Likes = append(pm.Posts[idx].Likes, newLikes...)
	}
	return nil
}

func (pm *PostsMockRepo) DeleteDislike(post entity.Post) error {
	newDislikes := []entity.Reaction{}
	idx := 0
	found := false
	for y, val := range pm.Posts {
		for i, v := range val.Dislikes {
			if v.UserId == post.User.Id {
				found = true
				newDislikes = deleteElement(val.Dislikes, i)
				idx = y
				break
			}
		}
	}

	if found {
		pm.Posts[idx].Dislikes = []entity.Reaction{}
		pm.Posts[idx].Dislikes = append(pm.Posts[idx].Dislikes, newDislikes...)
	}
	return nil
}

func (pm *PostsMockRepo) StoreTopicReference(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) GetRelatedCategories(post entity.Post) ([]string, error) {
	for _, v := range pm.Posts {
		if v.Id == post.Id {
			return v.Categories, nil
		}
	}
	return nil, errNoRows
}

func (pm *PostsMockRepo) FetchReactions(id int64) (entity.Post, error) {
	for _, v := range pm.Posts {
		if v.Id == id {
			return v, nil
		}
	}
	return entity.Post{}, errNoRows
}

func (pm *PostsMockRepo) StoreCategories(categories []string) error {
	for _, v := range categories {
		pm.AllTopics[v] = true
	}
	return nil
}

func (pm *PostsMockRepo) GetExistedCategories() ([]string, error) {
	categories := []string{}
	for key := range pm.AllTopics {
		categories = append(categories, key)
	}
	return categories, nil
}

type CommentsMockRepo struct {
	Posts    []entity.Post
	Comments []entity.Comment
}

func NewCommentsMockrepo() *CommentsMockRepo {
	return &CommentsMockRepo{}
}

func (cm *CommentsMockRepo) Store(comment entity.Comment) error {
	cm.Comments = append(cm.Comments, comment)
	return nil
}

func (cm *CommentsMockRepo) Fetch(postId int64) ([]entity.Comment, error) {
	return cm.Comments, nil
}

func (cm *CommentsMockRepo) GetById(id int64) (entity.Comment, error) {
	for i := 0; i < len(cm.Comments); i++ {
		if cm.Comments[i].Id == id {
			return cm.Comments[i], nil
		}
	}
	return entity.Comment{}, errNoRows
}

func (cm *CommentsMockRepo) GetPostIds(user entity.User) ([]int64, error) {
	var ids []int64
	for _, val := range cm.Posts {
		for _, v := range val.Comments {
			if v.User.Id == user.Id {
				ids = append(ids, val.Id)
				break
			}
		}
	}
	return ids, nil
}

func (cm *CommentsMockRepo) Update(comment entity.Comment) error {
	for i := 0; i < len(cm.Comments); i++ {
		if cm.Comments[i].Id == comment.Id {
			cm.Comments[i] = comment
			return nil
		}
	}
	return errNoRows
}

func (cm *CommentsMockRepo) Delete(comment entity.Comment) error {
	newComments := []entity.Comment{}
	found := false
	for i, v := range cm.Comments {
		if v.Id == comment.Id {
			newComments = deleteElement(cm.Comments, i)
			found = true
		}
	}
	if found {
		cm.Comments = []entity.Comment{}
		cm.Comments = append(cm.Comments, newComments...)
		return nil
	} else {
		return errNoRows
	}
}

func (cm *CommentsMockRepo) StoreLike(comment entity.Comment) error {
	found := false
	like := entity.Reaction{UserId: comment.User.Id}
	for i := 0; i < len(cm.Comments); i++ {
		if cm.Comments[i].Id == comment.Id {
			alreadyReacted := false
			for _, v := range cm.Comments[i].Likes {
				if v.UserId == comment.User.Id {
					alreadyReacted = true
				}
			}
			if alreadyReacted {
				return errUniqueConstraint
			}
			cm.Comments[i].Likes = append(cm.Comments[i].Likes, like)
			found = true
		}
	}
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (cm *CommentsMockRepo) DeleteLike(comment entity.Comment) error {
	newLikes := []entity.Reaction{}
	Idx := 0
	found := false
	for y, val := range cm.Comments {
		for i, v := range val.Likes {
			if v.UserId == comment.User.Id {
				found = true
				newLikes = deleteElement(val.Likes, i)
				Idx = y
				break
			}
		}
	}

	if found {
		cm.Comments[Idx].Likes = []entity.Reaction{}
		cm.Comments[Idx].Likes = append(cm.Comments[Idx].Likes, newLikes...)
	}
	return nil
}

func (cm *CommentsMockRepo) StoreDislike(comment entity.Comment) error {
	found := false
	dislike := entity.Reaction{UserId: comment.User.Id}
	for i := 0; i < len(cm.Comments); i++ {
		if cm.Comments[i].Id == comment.Id {
			alreadyReacted := false
			for _, v := range cm.Comments[i].Dislikes {
				if v.UserId == comment.User.Id {
					alreadyReacted = true
				}
			}
			if alreadyReacted {
				return errUniqueConstraint
			}
			cm.Comments[i].Dislikes = append(cm.Comments[i].Dislikes, dislike)
			found = true
		}
	}
	if found {
		return nil
	} else {
		return errNoRows
	}
}

func (cm *CommentsMockRepo) DeleteDislike(comment entity.Comment) error {
	newDislikes := []entity.Reaction{}
	Idx := 0
	found := false
	for y, val := range cm.Comments {
		for i, v := range val.Dislikes {
			if v.UserId == comment.User.Id {
				found = true
				newDislikes = deleteElement(val.Dislikes, i)
				Idx = y
				break
			}
		}
	}

	if found {
		cm.Comments[Idx].Dislikes = []entity.Reaction{}
		cm.Comments[Idx].Dislikes = append(cm.Comments[Idx].Dislikes, newDislikes...)
	}
	return nil
}

func (cm *CommentsMockRepo) FetchReactions(id int64) (entity.Comment, error) {
	for _, v := range cm.Comments {
		if v.Id == id {
			return v, nil
		}
	}
	return entity.Comment{}, errNoRows
}

func deleteElement[C any](sl []C, index int) []C {
	newSlice := []C{}
	newSlice = append(newSlice, sl[:index]...)
	newSlice = append(newSlice, sl[index+1:]...)
	return newSlice
}
