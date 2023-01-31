package entity

import "time"

type User struct {
	Id              int64
	Name            string
	Email           string
	Password        string
	RegDate         string
	DateOfBirth     string
	City            string
	Owner           bool
	Gender          string
	Male            bool
	Female          bool
	Role            string
	Sign            string
	SessionToken    string
	SessionTTL      time.Time
	Posts           int64
	Comments        int64
	PostLikes       int64
	PostDislikes    int64
	CommentLikes    int64
	CommentDislikes int64
}
