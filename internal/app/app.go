package app

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
	"forum/pkg/sqlite3"
	"log"
)

func Run() {
	sq, err := sqlite3.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sq.Close()
	repositories := repository.NewRepositories(sq)
	err = sqlite.CreateDB(sq)
	if err != nil {
		log.Fatal(err)
	}
	hasher := hasher.NewBcryptHasher()
	tokenManager, err := auth.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	useCases := usecase.NewUseCases(usecase.Dependencies{
		Repos:        repositories,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	// regDate := "2022-11-10"
	// date, err := time.Parse("2006-01-02", regDate)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// dateOfBirth := "1989-06-19"
	// birthDate, err := time.Parse("2006-01-02", dateOfBirth)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// user := entity.User{
	// 	Name:        "Bobik",
	// 	Email:       "bobik@gmail.com",
	// 	Password:    "vivse",
	// 	RegDate:     date,
	// 	DateOfBirth: birthDate,
	// 	City:        "Astana",
	// 	Sex:         "Male",
	// }
	id := int64(4)
	userFind := entity.User{
		Id: id,
	}
	// err = useCases.Users.SignUp(user)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = useCases.Users.UpdateSession(userFind)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	expired, err := useCases.Users.CheckTTLExpired(userFind)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(expired)
	// err = useCases.Users.DeleteUser(userFind)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// findUser, err := useCases.Users.GetById(2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = useCases.Users.UpdateSession(findUser)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = useCases.Users.UpdateSession(user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// post := entity.Post{
	// 	User: entity.User{
	// 		Id: 1,
	// 	},
	// 	Date:    "05.11.2022",
	// 	Title:   "Cars",
	// 	Content: "sdfsdf",
	// 	Category: []string{
	// 		"cars",
	// 		"mars",
	// 	},
	// }
	// err = useCases.Posts.CreatePost(post)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fetchedUsers, err := useCases.Users.GetAllUsers()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i := 0; i < len(fetchedUsers); i++ {
	// 	fmt.Println(fetchedUsers[i])
	// }

	// userByid, err := useCases.Users.GetById(2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(userByid.Name)

	// UserSession, err := useCases.Users.GetSession(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(UserSession)
}
