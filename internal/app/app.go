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
	"time"
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
	// 	// Name:        "Bobik1",
	// 	Email:       "bobik@gmail.com",
	// 	Password:    "vivse",
	// 	RegDate:     date,
	// 	DateOfBirth: birthDate,
	// 	City:        "Astana",
	// 	Sex:         "Male",
	// }
	// id := int64(4)
	// userFind := entity.User{
	// 	Id: id,
	// }
	// err = useCases.Users.SignUp(user)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = useCases.Users.SignIn(user)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = useCases.Users.UpdateSession(userFind)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// expired, err := useCases.Users.CheckTTLExpired(userFind)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(expired)
	// err = useCases.Users.DeleteUser(userFind)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// userP := entity.User{
	// 	Id: 1,
	// }
	// date := "2022-11-10 15:00:01"
	// parsed, err := time.Parse("2006-01-02 15:04:05", date)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// post := entity.Post{
	// 	User:    userP,
	// 	Date:    parsed,
	// 	Title:   "Carsd",
	// 	Content: "sdfsddf",
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

	user := entity.User{
		Id: 1,
	}
	date := "2022-11-07 15:00:45"
	parsed, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		log.Fatal(err)
	}
	post := entity.Post{
		Id:      1,
		Title:   "new title",
		Content: "new content",
		User:    user,
		Date:    parsed,
	}
	err = useCases.Posts.DeletePost(post)
	if err != nil {
		log.Fatal(err)
	}
	// post, err := useCases.Posts.GetById(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(post.Id, post.Title, post.Content, post.Date, post.User.Id, post.User.Name)
}
