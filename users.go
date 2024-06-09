package main

import "errors"

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
		Id:    newUser.Id,
		Email: newUser.Email,
	}

	return userOut, err
}

func (db *DB) GetUser(email string) (User, error) {
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

func (db *DB) UserExists(email string) bool {
	_, err := db.GetUser(email)
	return err == nil
}
