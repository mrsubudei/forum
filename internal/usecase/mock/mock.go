package mock_usecase

import (
	"forum/internal/entity"
	"time"
)

type UsersMockUseCase struct {
	Users []entity.User
}

func NewUsersMockUseCase() *UsersMockUseCase {
	return &UsersMockUseCase{}
}

func (um *UsersMockUseCase) SignUp(u entity.User) error {
	um.Users = append(um.Users, u)
	return nil
}

func (um *UsersMockUseCase) SignIn(u entity.User) error {
	return nil
}

func (um *UsersMockUseCase) GetAllUsers() ([]entity.User, error) {
	return []entity.User{}, nil
}

func (um *UsersMockUseCase) GetById(id int64) (entity.User, error) {
	return entity.User{}, nil
}

func (um *UsersMockUseCase) GetNewToken() (string, time.Time, error) {
	return "", time.Time{}, nil
}

func (um *UsersMockUseCase) GetIdBy(user entity.User) (int64, error) {
	var id int64
	if len(um.Users) > 1 {
		id = 5
	} else if len(um.Users) == 1 {
		id = 1
	}

	return id, nil
}

func (um *UsersMockUseCase) GetSession(id int64) (entity.User, error) {
	return entity.User{}, nil
}

func (um *UsersMockUseCase) CheckSession(u entity.User) (bool, error) {
	if len(um.Users) != 0 {
		return true, nil
	}
	return false, nil
}

func (um *UsersMockUseCase) UpdateUserInfo(u entity.User, query string) error {
	return nil
}

func (um *UsersMockUseCase) UpdateSession(u entity.User) error {
	return nil
}

func (um *UsersMockUseCase) DeleteSession(user entity.User) error {
	return nil
}

func (um *UsersMockUseCase) DeleteUser(u entity.User) error {
	um.Users = um.Users[:len(um.Users)-1]
	return nil
}

type PostsMockUseCase struct {
	Categories []string
}

func NewPostsMockUseCase() *PostsMockUseCase {
	return &PostsMockUseCase{}
}

func (pm *PostsMockUseCase) CreatePost(p entity.Post) error {
	return nil
}

func (pm *PostsMockUseCase) GetAllPosts() ([]entity.Post, error) {
	return []entity.Post{}, nil
}

func (pm *PostsMockUseCase) GetPostsByQuery(user entity.User, query string) ([]entity.Post, error) {
	return []entity.Post{}, nil
}

func (pm *PostsMockUseCase) GetById(id int64) (entity.Post, error) {
	return entity.Post{}, nil
}

func (pm *PostsMockUseCase) GetAllByCategory(category string) ([]entity.Post, error) {
	posts := []entity.Post{}
	found := false

	for _, v := range pm.Categories {
		if v == category {
			found = true
			break
		}
	}

	if !found {
		return posts, entity.ErrPostNotFound
	}
	return posts, nil
}

func (pm *PostsMockUseCase) UpdatePost(post entity.Post) error {
	return nil
}

func (pm *PostsMockUseCase) DeletePost(p entity.Post) error {
	return nil
}

func (pm *PostsMockUseCase) MakeReaction(p entity.Post, command string) error {
	return nil
}

func (pm *PostsMockUseCase) DeleteReaction(post entity.Post, command string) error {
	return nil
}

func (pm *PostsMockUseCase) CreateCategories(categories []string) error {
	pm.Categories = append(pm.Categories, categories...)
	return nil
}

func (pm *PostsMockUseCase) GetAllCategories() ([]string, error) {
	return []string{}, nil
}

func (pm *PostsMockUseCase) GetReactions(id int64, query string) ([]entity.User, error) {
	return []entity.User{}, nil
}

type CommentsMockUseCase struct{}

func NewCommentsMockUseCase() *CommentsMockUseCase {
	return &CommentsMockUseCase{}
}

func (cm *CommentsMockUseCase) WriteComment(c entity.Comment) error {
	return nil
}

func (cm *CommentsMockUseCase) GetAllComments(postId int64) ([]entity.Comment, error) {
	return []entity.Comment{}, nil
}

func (cm *CommentsMockUseCase) UpdateComment(c entity.Comment) error {
	return nil
}

func (cm *CommentsMockUseCase) DeleteComment(c entity.Comment) error {
	return nil
}

func (cm *CommentsMockUseCase) MakeReaction(c entity.Comment, command string) error {
	return nil
}

func (cm *CommentsMockUseCase) DeleteReaction(c entity.Comment, command string) error {
	return nil
}

func (cm *CommentsMockUseCase) GetReactions(id int64, query string) ([]entity.User, error) {
	return []entity.User{}, nil
}
