package models

import "time"

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	Profile        Profile
}

type Profile struct {
	UserId    int
	Avatar    string
	CreatedAt time.Time
}
