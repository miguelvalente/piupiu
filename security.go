package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(user User, expire int) string {
	var expiresAt time.Time

	if expire == 0 {
		expiresAt = time.Now().Add(24 * time.Hour)
	} else {
		expiresAt = time.Now().Add(time.Duration(expire) * time.Second)
	}
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
