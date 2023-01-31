package entity

type Post struct {
	Id               int64
	User             User
	Date             string
	Title            string
	Content          string
	ImagePath        string
	ContentWeb       []string
	Categories       []string
	Comments         []Comment
	LastComment      Comment
	LastCommentExist bool
	TotalComments    int64
	Likes            []Reaction
	TotalLikes       int64
	Dislikes         []Reaction
	TotalDislikes    int64
}

type Reaction struct {
	UserId int64  `json:"user_id,omitempty"`
	Date   string `json:"reaction_date,omitempty"`
}
