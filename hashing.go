package main

import (
	"golang.org/x/crypto/bcrypt"
)

func hash(clear_password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(clear_password), 10)
	if err != nil {
		return "", err
	}

	return string(hashed_password), nil
}
