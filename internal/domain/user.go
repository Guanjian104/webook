package domain

import "time"

type User struct {
    Id          int64
    Email       string
    Password    string
    Nickname    string
    Birthday    string
    Description string
    Phone       string

    Ctime time.Time
}

type UserProfile struct {
    Nickname    string
    Birthday    string
    Description string
}
