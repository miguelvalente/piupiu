package main

import "time"

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}
type User struct {
	Id             int       `json:"id"`
	Email          string    `json:"email"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
	Password       string    `json:"password"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresRefresh time.Time `json:"expires_in_seconds_refresh,omitempty"`
	Expires        int       `json:"expires_in_seconds,omitempty"`
}

type UserOut struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserOutLogin struct {
	Email        string `json:"email"`
	UserId       int    `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}
