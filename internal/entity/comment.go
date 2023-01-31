package entity

type Comment struct {
	Id            int64
	PostId        int64
	User          User
	Date          string
	Content       string
	ImagePath     string
	ContentWeb    []string
	Likes         []Reaction
	TotalLikes    int64
	Dislikes      []Reaction
	TotalDislikes int64
}
