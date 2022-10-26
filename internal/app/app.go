package app

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/internal/usecase"
	"forum/internal/usecase/repo/sqlite"
	"forum/pkg/sqlite3"
	"log"
	"os"
)

func Run() {
	sq, err := sqlite3.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sq.Close()
	CommunicationUseCase := usecase.New(sqlite.New(sq))

	err = CommunicationUseCase.CreateDB()
	if err != nil {
		log.Fatal(err)
	}

	topics := []string{"cars", "sports", "cinema"}
	CommunicationUseCase.CreateTopics(topics)

	user := entity.User{
		Name:        "Sub1i",
		Email:       "Sub1udei@gmail.com",
		Password:    "vivse",
		RegDate:     "10.10.2022",
		DateOfBirth: "19.06.1989",
		City:        "Astana",
		Sex:         "Male",
	}

	assignedUserId, err := CommunicationUseCase.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
	}
	user.Id = assignedUserId

	user2 := entity.User{
		Name:        "Mu11bi",
		Email:       "Mub1udei@gmail.com",
		Password:    "vivse",
		RegDate:     "10.10.2022",
		DateOfBirth: "19.06.1989",
		City:        "Astana",
		Sex:         "Male",
	}

	assignedUserId, err = CommunicationUseCase.CreateUser(&user2)
	if err != nil {
		fmt.Println(err)
	}

	user2.Id = assignedUserId
	post := entity.Post{
		UserId:  1,
		Date:    "25.10.2022",
		Content: "Drop file here to load content or click on this box to open file dialog.",
	}

	assignedPostId, err := CommunicationUseCase.CreatePost(&post)
	if err != nil {
		fmt.Println(err)
	}

	post.Id = assignedPostId

	post2 := entity.Post{
		UserId:  1,
		Date:    "25.10.2022",
		Content: "Droptfhdtyhftghjlog.",
	}

	assignedPostId, err = CommunicationUseCase.CreatePost(&post2)
	if err != nil {
		fmt.Println(err)
	}

	post2.Id = assignedPostId

	err = CommunicationUseCase.CreatePostRef(&post, []string{"cars", "sports"})
	if err != nil {
		fmt.Println(err)
	}

	err = CommunicationUseCase.CreateComment(&post, "01.01.2022", "good")
	if err != nil {
		fmt.Println(err)
	}

	err = CommunicationUseCase.PutPostLike(&user, 1, "25.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = CommunicationUseCase.PutPostDisLike(&user, 1, "24.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = CommunicationUseCase.PutCommentLike(&user, 1, "20.10.2022")
	if err != nil {
		fmt.Println(err)
	}

	err = CommunicationUseCase.PutCommentDisLike(&user, 1, "21.10.2022")
	if err != nil {
		fmt.Println(err)
	}
	users, err := CommunicationUseCase.GetAllUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Printint all users:")
	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
	}
	fmt.Println("Printint selected user:")
	selectedUser, err := (CommunicationUseCase.GetUser(1))
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
