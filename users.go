package main

import (
	"errors"
	"fmt"
	"time"
)

func (db *DB) CreateUser(email string, hashed_password string) (UserOut, error) {
	err := db.ensureDB()
	if err != nil {
		return UserOut{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserOut{}, err
	}

	userId := len(dbStructure.Users) + 1

	newUser := User{
		Id:       userId,
		Email:    email,
		Password: hashed_password,
	}
	if dbStructure.Users == nil {
		dbStructure.Users = map[int]User{}
	}

	dbStructure.Users[userId] = newUser
	err = db.writeDB(dbStructure)

	userOut := UserOut{
		Id:    userId,
		Email: newUser.Email,
	}

	return userOut, err
}

func (db *DB) UpdateUser(user_id int, email string, hashed_password string) (UserOut, error) {
	err := db.ensureDB()
	if err != nil {
		return UserOut{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserOut{}, err
	}

	userExisting, err := db.GetUserById(user_id)

	dbStructure.Users[user_id] = User{
		Id:             user_id,
		Email:          email,
		Password:       hashed_password,
		RefreshToken:   userExisting.RefreshToken,
		ExpiresRefresh: userExisting.ExpiresRefresh,
	}

	err = db.writeDB(dbStructure)

	userOut := UserOut{
		Id:    user_id,
		Email: email,
	}

	return userOut, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	err := db.ensureDB()
	if err != nil {
		return User{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("User email not found")

}

func (db *DB) GetUserById(id int) (User, error) {
	err := db.ensureDB()
	if err != nil {
		return User{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if user, exists := dbStructure.Users[id]; exists {
		return user, nil
	} else {
		return User{}, errors.New("User email not found")
	}

}

func (db *DB) AddRefreshTokenToUser(user User, token string) error {
	err := db.ensureDB()
	if err != nil {
		return err
	}

	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	var expiresAt time.Time
	expiresAt = time.Now().Add(1440 * time.Hour).UTC()

	user.RefreshToken = token
	user.ExpiresRefresh = expiresAt

	dbStructure.Users[user.Id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (db *DB) UserExists(email string) bool {
	_, err := db.GetUserByEmail(email)
	return err == nil
}

func (db *DB) RefreshTokenValid(token string) (string, bool, error) {
	db.ensureDB()

	dbStructure, err := db.loadDB()
	if err != nil {
		return "", false, err
	}

	nowTime := time.Now().UTC()

	for _, user := range dbStructure.Users {
		if user.RefreshToken == token {
			fmt.Println("0", user.ExpiresRefresh)
			fmt.Println("1", nowTime)
			if user.ExpiresRefresh.After(nowTime) {
				newToken := createJWT(user)
				return newToken, true, nil
			}
		}
	}

	return "", false, err
}

func (db *DB) RevokeToken(token string) error {
	db.ensureDB()

	dbStructure, err := db.loadDB()
	// fmt.Println("befiore")
	// fmt.Println(dbStructure)
	if err != nil {
		return err
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken == token {
			if user.ExpiresRefresh.After(time.Now()) {
				user.ExpiresRefresh = time.Now().AddDate(0, -1, 0).UTC()
				dbStructure.Users[user.Id] = user
			}
		}
	}
	// fmt.Println("after")
	// fmt.Println(dbStructure)
	db.writeDB(dbStructure)
	dbStructure, err = db.loadDB()
	// fmt.Println("after after")
	// fmt.Println(dbStructure)

	return nil
}
