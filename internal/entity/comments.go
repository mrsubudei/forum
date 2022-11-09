package entity

import "time"

type Comment struct {
	Id            int64      `json:"id,omitempty"`
	PostId        int64      `json:"post_id,omitempty"`
	UserId        int64      `json:"user_id,omitempty"`
	UserName      string     `json:"user_name,omitempty"`
	Date          time.Time  `json:"comment_date,omitempty"`
	Content       string     `json:"comment_content,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	TotalLikes    int64      `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislikes,omitempty"`
	TotalDislikes int64      `json:"total_dislikes,omitempty"`
}
