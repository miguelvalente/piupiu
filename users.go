package main

func (db *DB) CreateUser(email string) (User, error) {
	err := db.ensureDB()
	if err != nil {
		return User{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	userId := len(dbStructure.Users) + 1

	newUser := User{
		Id:    userId,
		Email: email,
	}
	if dbStructure.Users == nil {
		dbStructure.Users = map[int]User{}
	}

	dbStructure.Users[userId] = newUser
	err = db.writeDB(dbStructure)

	return newUser, err
}
