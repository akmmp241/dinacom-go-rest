package model

import "time"

type User struct {
	Id       int
	Email    string
	Password string
}

type Session struct {
	Id        int
	UserId    int
	Token     string
	ExpiresAt time.Time
}

type Complaint struct {
	Id            string
	UserId        int
	Title         string
	ComplaintsMsg string
	Response      string
	ImageUrl      string
	CreatedAt     time.Time
}
