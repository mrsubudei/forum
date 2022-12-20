package mock_usecase

import (
	"forum/internal/entity"
	"forum/internal/usecase"
)

type MockUseCases struct {
	Users    usecase.Users
	Posts    usecase.Posts
	Comments usecase.Comments
}

func NewMockUseCases() *MockUseCases {
	return &MockUseCases{
		Users:    NewUsersMockUseCase(),
		Posts:    NewPostsMockUseCase(),
		Comments: NewCommentsMockUseCase(),
	}
}

type UsersMockUseCase struct{}

func NewUsersMockUseCase() *UsersMockUseCase {
	return &UsersMockUseCase{}
}

func (um *UsersMockUseCase) SignUp(u entity.User) error {
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

func (um *UsersMockUseCase) GetIdBy(user entity.User) (int64, error) {
	var id int64
	return id, nil
}

func (um *UsersMockUseCase) GetSession(id int64) (entity.User, error) {
	return entity.User{}, nil
}

func (um *UsersMockUseCase) CheckSession(u entity.User) (bool, error) {
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
	return nil
}

type PostsMockUseCase struct{}

func NewPostsMockUseCase() *PostsMockUseCase {
	return &PostsMockUseCase{}
}

func (pm *PostsMockUseCase) CreatePost(p entity.Post) error
func (pm *PostsMockUseCase) GetAllPosts() ([]entity.Post, error)
func (pm *PostsMockUseCase) GetPostsByQuery(user entity.User, query string) ([]entity.Post, error)
func (pm *PostsMockUseCase) GetById(id int64) (entity.Post, error)
func (pm *PostsMockUseCase) GetAllByCategory(category string) ([]entity.Post, error)
func (pm *PostsMockUseCase) UpdatePost(post entity.Post) error
func (pm *PostsMockUseCase) DeletePost(p entity.Post) error
func (pm *PostsMockUseCase) MakeReaction(p entity.Post, command string) error
func (pm *PostsMockUseCase) DeleteReaction(post entity.Post, command string) error
func (pm *PostsMockUseCase) CreateCategories(categories []string) error
func (pm *PostsMockUseCase) GetAllCategories() ([]string, error)
func (pm *PostsMockUseCase) GetReactions(id int64, query string) ([]entity.User, error)

type CommentsMockUseCase struct{}

func NewCommentsMockUseCase() *CommentsMockUseCase {
	return &CommentsMockUseCase{}
}

func (cm *CommentsMockUseCase) WriteComment(c entity.Comment) error
func (cm *CommentsMockUseCase) GetAllComments(postId int64) ([]entity.Comment, error)
func (cm *CommentsMockUseCase) UpdateComment(c entity.Comment) error
func (cm *CommentsMockUseCase) DeleteComment(c entity.Comment) error
func (cm *CommentsMockUseCase) MakeReaction(c entity.Comment, command string) error
func (cm *CommentsMockUseCase) DeleteReaction(c entity.Comment, command string) error
func (cm *CommentsMockUseCase) GetReactions(id int64, query string) ([]entity.User, error)
