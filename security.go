package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(user User) string {
	var expiresAt time.Time
	expiresAt = time.Now().Add(1 * time.Hour)

	// if expire == 0 {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   fmt.Sprint(user.Id),
		ExpiresAt: jwt.NewNumericDate(expiresAt.UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")

	string, _ := token.SignedString([]byte(jwtSecret))

	return string
}
