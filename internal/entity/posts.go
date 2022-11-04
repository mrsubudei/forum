package entity

type Post struct {
	Id            int        `json:"id,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"date,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Category      []Category `json:"categories,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	CountComments int        `json:"total_comments,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	CountLikes    int        `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislkes,omitempty"`
	CountDislikes int        `json:"total_dislikes,omitempty"`
}

type Reaction struct {
	Date string `json:"reaction_date,omitempty"`
	User User   `json:"user,omitempty"`
}

type Category struct {
	Title      string `json:"title,omitempty"`
	Posts      []Post `json:"posts,omitempty"`
	CountPosts int    `json:"total_posts,omitempty"`
}
