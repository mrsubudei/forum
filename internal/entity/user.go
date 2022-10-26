package entity

type User struct {
	Id          int
	Name        string
	Email       string
	Password    string
	RegDate     string
	DateOfBirth string
	City        string
	Sex         string
	PostIds     []int
}
