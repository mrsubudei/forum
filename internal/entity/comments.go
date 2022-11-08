package entity

import "time"

type Comment struct {
	Id            int64      `json:"id,omitempty"`
	Post          Post       `json:"post,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          time.Time  `json:"comment_date,omitempty"`
	Content       string     `json:"comment_content,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	TotalLikes    int64      `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislikes,omitempty"`
	TotalDislikes int64      `json:"total_dislikes,omitempty"`
}
