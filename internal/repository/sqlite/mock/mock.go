package mock_repository

import (
	"fmt"

	"forum/internal/entity"
	"forum/internal/usecase"
)

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
	users := []entity.User{}
	return users, nil
}

func (um *UsersMockRepo) GetId(user entity.User) (int64, error) {
	var id int64
	return id, nil
}

func (um *UsersMockRepo) GetById(n int64) (entity.User, error) {
	user := entity.User{}
	return user, nil
}

func (um *UsersMockRepo) GetSession(n int64) (entity.User, error) {
	user := entity.User{}
	return user, nil
}

func (um *UsersMockRepo) UpdateInfo(user entity.User) error {
	return nil
}

func (um *UsersMockRepo) UpdatePassword(user entity.User) error {
	return nil
}

func (um *UsersMockRepo) NewSession(user entity.User) error {
	return nil
}

func (um *UsersMockRepo) UpdateSession(user entity.User) error {
	return nil
}

func (um *UsersMockRepo) Delete(user entity.User) error {
	return nil
}

type PostsMockRepo struct {
	posts []entity.Post
}

func NewPostsMockrepo() *PostsMockRepo {
	return &PostsMockRepo{}
}

func (pm *PostsMockRepo) Store(post *entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) Fetch() ([]entity.Post, error) {
	posts := []entity.Post{}
	return posts, nil
}

func (pm *PostsMockRepo) FetchByAuthor(user entity.User) ([]entity.Post, error) {
	posts := []entity.Post{}
	return posts, nil
}

func (pm *PostsMockRepo) GetById(id int64) (entity.Post, error) {
	post := entity.Post{}
	return post, nil
}

func (pm *PostsMockRepo) GetIdsByCategory(category string) ([]int64, error) {
	var ids []int64
	return ids, nil
}

func (pm *PostsMockRepo) FetchIdsByReaction(user entity.User, reaction string) ([]int64, error) {
	var ids []int64
	return ids, nil
}

func (pm *PostsMockRepo) Update(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) Delete(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) StoreLike(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) StoreDislike(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) DeleteLike(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) DeleteDislike(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) StoreTopicReference(post entity.Post) error {
	return nil
}

func (pm *PostsMockRepo) GetRelatedCategories(post entity.Post) ([]string, error) {
	var categories []string
	return categories, nil
}

func (pm *PostsMockRepo) FetchReactions(id int64) (entity.Post, error) {
	post := entity.Post{}
	return post, nil
}

func (pm *PostsMockRepo) StoreCategories(categories []string) error {
	return nil
}

func (pm *PostsMockRepo) GetExistedCategories() ([]string, error) {
	var categories []string
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
