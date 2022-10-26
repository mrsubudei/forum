package usecase

import (
	"forum/internal/entity"
)

type CommunicationUseCase struct {
	repo CommunicationRepo
}

func New(r CommunicationRepo) *CommunicationUseCase {
	return &CommunicationUseCase{
		repo: r,
	}
}

func (c *CommunicationUseCase) CreateDB() error {
	err := c.repo.CreateDB()
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) GetAllUsers() ([]entity.User, error) {
	users, err := c.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *CommunicationUseCase) GetUser(userId int) (entity.User, error) {
	user, err := c.repo.GetUser(userId)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (c *CommunicationUseCase) GetUserPostIds(userId int) ([]int, error) {
	ids, err := c.repo.GetUserPostIds(userId)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (c *CommunicationUseCase) GetAllPosts() ([]entity.Post, error) {
	posts, err := c.repo.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *CommunicationUseCase) CreateUser(u *entity.User) (int, error) {
	userId, err := c.repo.CreateUser(u)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (c *CommunicationUseCase) CreatePost(p *entity.Post) (int, error) {
	postId, err := c.repo.CreatePost(p)
	if err != nil {
		return 0, err
	}
	return postId, nil
}

func (c *CommunicationUseCase) CreateComment(p *entity.Post, date, content string) error {
	err := c.repo.CreateComment(p, date, content)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) PutPostLike(u *entity.User, postId int, date string) error {
	err := c.repo.PutPostLike(u, postId, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) PutPostDisLike(u *entity.User, postId int, date string) error {
	err := c.repo.PutPostDisLike(u, postId, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) PutCommentLike(u *entity.User, commentId int, date string) error {
	err := c.repo.PutCommentLike(u, commentId, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) PutCommentDisLike(u *entity.User, commentId int, date string) error {
	err := c.repo.PutCommentDisLike(u, commentId, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) CreateTopics(name []string) error {
	err := c.repo.CreateTopics(name)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunicationUseCase) CreatePostRef(p *entity.Post, name []string) error {
	err := c.repo.CreatePostRef(p, name)
	if err != nil {
		return err
	}
	return nil
}
