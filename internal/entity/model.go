package entity

type Post struct {
	Id       int        `json:"post_id"`
	User     User       `json:"user"`
	Date     string     `json:"post_date"`
	Content  string     `json:"post_content"`
	Topics   []string   `json:"topics"`
	Comments []Comment  `json:"comments"`
	Likes    []Reaction `json:"likes"`
	Dislikes []Reaction `json:"dislkes"`
}

type Comment struct {
	Id       int        `json:"comment_id"`
	User     User       `json:"user"`
	Date     string     `json:"comment_date"`
	Content  string     `json:"comment_content"`
	Likes    []Reaction `json:"likes"`
	Dislikes []Reaction `json:"dislikes"`
}

type Reaction struct {
	Like bool   `json:"like"`
	Date string `json:"reaction_date"`
	User User   `json:"user"`
}

type PostReaction struct {
	Reaction `json:"reaction"`

	Post `json:"post"`
}

type CommentReaction struct {
	Reaction `json:"reaction"`
	Comment  `json:"comment"`
}

type User struct {
	Id              int               `json:"user_id"`
	Name            string            `json:"user_name"`
	Email           string            `json:"email"`
	Password        string            `json:"password"`
	RegDate         string            `json:"registration_date"`
	DateOfBirth     string            `json:"date_of_birth"`
	City            string            `json:"city"`
	Sex             string            `json:"sex"`
	Posts           []Post            `json:"posts"`
	PostReactions   []PostReaction    `json:"post_reactions"`
	CommentReaction []CommentReaction `json:"comment_reactions"`
}
