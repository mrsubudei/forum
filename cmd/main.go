package main

import (
	"errors"
	"fmt"
	u "forum/internal/database/usecases"
	"log"
	"os"
)

func main() {

	exist, err := Exists("forum.db")
	if err != nil {
		log.Fatal(err)
	}

	if !exist {
		u.CreateDB()
	}

	u.CreateTopics("cars")
	u.CreateTopics("sports")
	u.CreateTopics("cinema")

	user := u.User{
		Name:        "Subi",
		Email:       "Subudei@gmail.com",
		Password:    "vivse",
		RegDate:     "10.10.2022",
		DateOfBirth: "19.06.1989",
		City:        "Astana",
		Sex:         "Male",
	}

	assignedUserId, err := user.CreateUser()
	if err != nil {
		fmt.Println(err)
	}
	user.Id = assignedUserId

	user2 := u.User{
		Name:        "Mubi",
		Email:       "Mubudei@gmail.com",
		Password:    "vivse",
		RegDate:     "10.10.2022",
		DateOfBirth: "19.06.1989",
		City:        "Astana",
		Sex:         "Male",
	}

	assignedUserId, err = user2.CreateUser()
	if err != nil {
		fmt.Println(err)
	}

	user2.Id = assignedUserId
	post := u.Post{
		UserId:  1,
		Date:    "25.10.2022",
		Content: "Drop file here to load content or click on this box to open file dialog.",
	}

	assignedPostId, err := post.CreatePost()
	if err != nil {
		fmt.Println(err)
	}

	post.Id = assignedPostId

	post2 := u.Post{
		UserId:  1,
		Date:    "25.10.2022",
		Content: "Droptfhdtyhftghjlog.",
	}

	assignedPostId, err = post2.CreatePost()
	if err != nil {
		fmt.Println(err)
	}

	post2.Id = assignedPostId

	err = post.CreatePostRef("cars", "sports")
	if err != nil {
		fmt.Println(err)
	}

	err = post.CreateComment("01.01.2022", "good")
	if err != nil {
		fmt.Println(err)
	}

	err = user.PutPostLike(1, "25.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = user.PutPostDisLike(1, "24.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = user.PutCommentLike(1, "20.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = user.PutCommentDisLike(1, "21.10.2022")
	if err != nil {
		fmt.Println(err)
	}
	users, err := u.GetAllUsers()
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
	}
	selectedUser, err := (u.GetUser(1))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(selectedUser)

}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
