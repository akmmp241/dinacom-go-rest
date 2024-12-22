package model

import "time"

type User struct {
	Id       int
	Name     string
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
	CreatedAt     time.Time
}
