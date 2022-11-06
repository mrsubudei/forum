package entity

type Post struct {
	Id            int64      `json:"id,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"date,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Category      []string   `json:"categories,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	CountComments int64      `json:"total_comments,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	Dislikes      []Reaction `json:"dislkes,omitempty"`
}

type Reaction struct {
	UserId   int64  `json:"user_id,omitempty"`
	EntityId int64  `json:"entity,omitempty"`
	Date     string `json:"reaction_date,omitempty"`
}
