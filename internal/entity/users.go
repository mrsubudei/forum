package entity

type User struct {
	Id                    int               `json:"id,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Email                 string            `json:"email,omitempty"`
	Password              string            `json:"password,omitempty"`
	RegDate               string            `json:"registration_date,omitempty"`
	DateOfBirth           string            `json:"date_of_birth,omitempty"`
	City                  string            `json:"city,omitempty"`
	Sex                   string            `json:"sex,omitempty"`
	SessionToken          string            `json:"token,omitempty"`
	TokenExpiration       string            `json:"token_expiration,omitempty"`
	Posts                 []Post            `json:"posts,omitempty"`
	CountPosts            int               `json:"total_posts,omitempty"`
	Comments              []Comment         `json:"comments,omitempty"`
	CountComments         int               `json:"total_comments,omitempty"`
	PostReactions         []PostReaction    `json:"post_reactions,omitempty"`
	CountPostReactions    int               `json:"total_post_reactions,omitempty"`
	CommentReactions      []CommentReaction `json:"comment_reactions,omitempty"`
	CountCommentReactions int               `json:"total_comment_reactions,omitempty"`
}

type PostReaction struct {
	Reaction `json:"reaction,omitempty"`
	Post     `json:"post,omitempty"`
}

type CommentReaction struct {
	Reaction `json:"reaction,omitempty"`
	Comment  `json:"comment,omitempty"`
}
