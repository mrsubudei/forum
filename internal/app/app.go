package app

import (
	"fmt"
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
	categ := []string{"cars", "weather"}
	err = sqlite.CreateCategories(sq, categ)
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
	// dateOfBirth := "1989-01-19"
	// birthDate, err := time.Parse("2006-01-02", dateOfBirth)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// user := entity.User{
	// 	Name:        "Zhorik",
	// 	Email:       "Zhor@gmail1.com",
	// 	Password:    "vivsef",
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
	// 	Id: 2,
	// }
	// date := "2022-11-10 15:00:01"
	// parsed, err := time.Parse("2006-01-02 15:04:05", date)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// date := "2022-11-07 15:00:45"
	// parsed, err := time.Parse("2006-01-02 15:04:05", date)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// post := entity.Post{
	// 	Id: 1,
	// 	User: entity.User{
	// 		Id: 3,
	// 	},
	// 	Date:    parsed,
	// 	Title:   "Pogoda",
	// 	Content: "sdfsddfsdgdfg",
	// 	Categories: []string{
	// 		"cinema",
	// 		"weather",
	// 	},
	// }
	// comment := entity.Comment{
	// 	Id:     1,
	// 	UserId: 3,
	// 	Date:   parsed,
	// }
	// err = useCases.Comments.MakeReaction(comment, "like")
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

	// user := entity.User{
	// 	Id: 1,
	// }

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// post := entity.Post{
	// 	Id:   1,
	// 	User: user,
	// }
	// comment := entity.Comment{
	// 	PostId:  2,
	// 	UserId:  3,
	// 	Date:    parsed,
	// 	Content: "pshel nah",
	// }
	// err = useCases.Comments.WriteComment(comment)
	// posts, err := useCases.Posts.GetAllPosts()

	// comment := entity.Comment{
	// 	Id:      1,
	// 	Post:    post,
	// 	User:    user,
	// 	Date:    parsed,
	// 	Content: "sdfgdf",
	// }
	// err = useCases.Comments.MakeReaction(comment, "dislike")
	// err = useCases.Comments.WriteComment(comment)
	// post, err := useCases.Posts.GetById(2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for i := 0; i < len(posts); i++ {
	// 	fmt.Printf("%#v\n", posts[i])
	// }

	// fmt.Println("post id: ", post.Id)
	// fmt.Println("post user id: ", post.User.Id)
	// fmt.Println("post date: ", post.Date)
	// fmt.Println("post title: ", post.Title)
	// fmt.Println("post content: ", post.Content)
	// fmt.Println("post categories: ", post.Categories)
	// fmt.Println("post count comments: ", post.CountComments)
	// fmt.Println("post likes: ", post.TotalLikes)
	// fmt.Println("post dislikes: ", post.TotalDislikes)

	// for i := 0; i < len(post.Comments); i++ {
	// 	fmt.Printf("comment %d %#v\n", i, post.Comments[i])
	// }

	// comments, err := useCases.Comments.GetAllComments()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i := 0; i < len(comments); i++ {

	// 	fmt.Print(comments[i].TotalLikes)
	// 	fmt.Print(comments[i].TotalDislikes)
	// 	fmt.Println()
	// }
	// post, err := useCases.Posts.GetById(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(post.Id, post.Title, post.Content, post.Date, post.User.Id, post.User.Name)
	posts, err := useCases.Posts.GetOneByCategory("cahrs")
	if err != nil {
		log.Fatal(err)
	}
	// for i := 0; i < len(posts); i++ {
	fmt.Print("id:", posts.Id, " ")
	fmt.Print("user id:", posts.User.Id, " ")
	fmt.Print("user name:", posts.User.Name, " ")
	fmt.Print("title:", posts.Title, " ")
	fmt.Print("content", posts.Content, " ")
	fmt.Print("categories", posts.Categories, " ")
	fmt.Print("total comments", posts.TotalComments, " ")
	fmt.Print("total likes", posts.TotalLikes, " ")
	fmt.Print("total dislikes", posts.TotalDislikes, " ")
	fmt.Println()
	// }
}
