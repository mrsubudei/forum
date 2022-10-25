package usecases

type User struct {
	Id          int
	Name        string
	Email       string
	Password    string
	RegDate     string
	DateOfBirth string
	City        string
	Sex         string
	PostIds     []int
}

type Post struct {
	Id         int
	UserId     int
	Date       string
	Content    string
	Topics     []string
	CommentIds []int
	Likes      []Like
	DisLikes   []DisLike
}

type Comment struct {
	Id       int
	UserId   int
	Date     string
	Content  string
	Likes    []Like
	DisLikes []DisLike
}

type Like struct {
	Date   string
	UserId int
}

type DisLike struct {
	Date   string
	UserId int
}
