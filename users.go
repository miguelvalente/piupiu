package main

import (
	"errors"
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

	dbStructure.Users[user_id] = User{
		Id:       user_id,
		Email:    email,
		Password: hashed_password,
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

func (db *DB) UserExists(email string) bool {
	_, err := db.GetUserByEmail(email)
	return err == nil
}
