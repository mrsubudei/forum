package usecases

import (
	"database/sql"
	"log"
)

func GetAllUsers() ([]User, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listUsers := []User{}
	for rows.Next() {

		var id int
		var name string
		var email string
		var pass string
		var regDate string
		var dateOfBirth string
		var city string
		var sex string
		var posts []int
		err = rows.Scan(&id, &name, &email, &pass, &regDate, &dateOfBirth, &city, &sex)
		if err != nil {
			return nil, err
		}
		posts, err = GetUserPostIds(db, id)
		if err != nil {
			return nil, err
		}

		user := User{
			Id:          id,
			Name:        name,
			Password:    pass,
			RegDate:     regDate,
			DateOfBirth: dateOfBirth,
			City:        city,
			Sex:         sex,
			PostIds:     posts,
		}
		listUsers = append(listUsers, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return listUsers, nil
}

func GetUser(userId int) (User, error) {
	selectedUser := User{}
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return selectedUser, err
	}

	rows, err := db.Query(`SELECT id, name, email, password, reg_date, 
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
		selectedUser.PostIds, err = GetUserPostIds(db, userId)
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

func GetUserPostIds(db *sql.DB, userId int) ([]int, error) {

	rows, err := db.Query("SELECT id FROM posts WHERE user_id = ?", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var posts []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, id)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return posts, nil
}

func GetAllPosts() ([]Post, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listPosts := []Post{}
	for rows.Next() {

		var id int
		var userId int
		var date string
		var content string
		var topics []string
		var comments []int
		var likes []Like
		var dislikes []DisLike
		err = rows.Scan(&id, &userId, &date, &content)
		if err != nil {
			return nil, err
		}

		//topics = GetAllTopics()

		post := Post{
			Id:         id,
			UserId:     userId,
			Date:       date,
			Content:    content,
			Topics:     topics,
			CommentIds: comments,
			Likes:      likes,
			DisLikes:   dislikes,
		}
		listPosts = append(listPosts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return listPosts, nil
}
