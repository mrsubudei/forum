package entity

type Post struct {
	Id            int64      `json:"id,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"date,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Categories    []string   `json:"categories,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	LastCOmment   Comment    `json:"last_comment,omitempty"`
	TotalComments int64      `json:"total_comments,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	TotalLikes    int64      `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislkes,omitempty"`
	TotalDislikes int64      `json:"total_dislikes,omitempty"`
}

type Reaction struct {
	UserId int64  `json:"user_id,omitempty"`
	Date   string `json:"reaction_date,omitempty"`
}
