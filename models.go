package main

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
