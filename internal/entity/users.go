package entity

import "time"

type User struct {
	Id              int64     `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Email           string    `json:"email,omitempty"`
	Password        string    `json:"password,omitempty"`
	RegDate         string    `json:"registration_date,omitempty"`
	DateOfBirth     string    `json:"date_of_birth,omitempty"`
	City            string    `json:"city,omitempty"`
	Gender          string    `json:"sex,omitempty"`
	Role            string    `json:"role,omitempty"`
	SessionToken    string    `json:"token,omitempty"`
	SessionTTL      time.Time `json:"token_expiration,omitempty"`
	Posts           int64     `json:"posts,omitempty"`
	Comments        int64     `json:"comments,omitempty"`
	PostLikes       int64     `json:"post_likes,omitempty"`
	PostDislikes    int64     `json:"post_dislikes,omitempty"`
	CommentLikes    int64     `json:"comment_likes,omitempty"`
	CommentDislikes int64     `json:"comment_dislikes,omitempty"`
}
