package v1

import (
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) MainHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	_, _ = fmt.Fprintln(w, "ok")
	// if r.URL.Path != "/" {
	// 	http.Error(w, "404: Page is Not Found", http.StatusNotFound)
	// 	return
	// }
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// html, err := template.ParseFiles("templates/index.html")
	// if err != nil {
	// 	http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	// user := entity.User{
	// 	Name: "Ivan",
	// }
	// err = html.Execute(w, user)
	// if err != nil {
	// 	http.Error(w, "404: Not Found", 404)
	// 	return
	// }

	//-------------------------------
	// regDate := "2022-11-10"
	// date, err := time.Parse("2006-01-02", regDate)
	// if err != nil {
	// 	log.Fatal(err)
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

	// for i := 0; i < len(posts); i++ {
	// 	fmt.Print("id:", posts[i].Id, " ")
	// 	fmt.Print("user id:", posts[i].User.Id, " ")
	// 	fmt.Print("user name:", posts[i].User.Name, " ")
	// 	fmt.Print("title:", posts[i].Title, " ")
	// 	fmt.Print("content", posts[i].Content, " ")
	// 	fmt.Print("categories", posts[i].Categories, " ")
	// 	fmt.Print("total comments", posts[i].TotalComments, " ")
	// 	fmt.Print("total likes", posts[i].TotalLikes, " ")
	// 	fmt.Print("total dislikes", posts[i].TotalDislikes, " ")
	// 	fmt.Println()
	// }
}
