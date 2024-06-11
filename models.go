package main

import "time"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresRefresh time.Time `json:"expires_in_seconds_refresh,omitempty"`
	Expires        int       `json:"expires_in_seconds,omitempty"`
}

type UserOut struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserOutLogin struct {
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
