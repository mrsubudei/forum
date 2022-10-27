package sqlite

import (
	"forum/internal/entity"
)

func (c *CommunicationRepo) GetAllUsers() ([]entity.User, error) {
	rows, err := c.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listUsers := []entity.User{}
	for rows.Next() {

		user := entity.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegDate, &user.DateOfBirth, &user.City, &user.Sex)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		listUsers = append(listUsers, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return listUsers, nil
}

func (c *CommunicationRepo) GetUser(userId int) (entity.User, error) {
	selectedUser := entity.User{}

	rows, err := c.DB.Query(`SELECT id, name, email, password, reg_date, 
	date_of_birth, city, sex FROM users WHERE id = ?`, userId)
	if err != nil {
		return selectedUser, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&selectedUser.Id, &selectedUser.Name, &selectedUser.Email, &selectedUser.Password,
			&selectedUser.RegDate, &selectedUser.DateOfBirth, &selectedUser.City, &selectedUser.Sex)
		if err != nil {
			return selectedUser, err
		}
		if err != nil {
			return selectedUser, err
		}
	}
	err = rows.Err()
	if err != nil {
		return selectedUser, err
	}
	return selectedUser, err
}
func gettopics()

func (c *CommunicationRepo) GetUserPostIds(userId int) ([]int, error) {
	rows, err := c.DB.Query("SELECT id FROM posts WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, id)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *CommunicationRepo) GetAllPosts() ([]entity.Post, error) {
	rows, err := c.DB.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listPosts := []entity.Post{}
	for rows.Next() {

		post := entity.Post{}
		err = rows.Scan(&id, &date, &content)
		if err != nil {
			return nil, err
		}

		// topics = GetAllTopics()

		listPosts = append(listPosts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return listPosts, nil
}
