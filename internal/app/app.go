package app

import (
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
	tokenManager, err := auth.NewManager("s")
	if err != nil {
		log.Fatal(err)
	}
	useCases := usecase.NewUseCases(usecase.Dependencies{
		Repos:        repositories,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	// user := entity.User{
	// 	Name:     "dgd",
	// 	Email:    "1223dfgd3432c",
	// 	Password: "11fg34c",
	// 	RegDate:  "19.05.1966",
	// }
	// err = useCases.Users.SignUp(user)
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

	user := entity.User{
		Id: 1,
	}
	err = useCases.Users.DeleteUser(user)
	if err != nil {
		log.Fatal(err)
	}
}
