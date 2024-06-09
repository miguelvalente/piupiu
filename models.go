package main

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Expires  int    `json:"expires_in_seconds,omitempty"`
}

type UserOut struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserOutLogin struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
