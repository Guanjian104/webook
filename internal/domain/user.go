package domain

import "time"

type User struct {
	Id          int64
	Email       string
	Password    string
	Nickname    string
	Birthday    string
	Description string

	Ctime time.Time
}

type UserProfile struct {
	Id          int64
	Email       string
	Nickname    string
	Birthday    string
	Description string
}