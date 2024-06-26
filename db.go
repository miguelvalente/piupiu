package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		empty_db, _ := json.Marshal(new(DBStructure))

		err := os.WriteFile(path, empty_db, 0666)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	return &db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	err := db.ensureDB()
	if err != nil {
		return Chirp{}, err
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirpId := len(dbStructure.Chirps) + 1

	newChirp := Chirp{
		Id:       chirpId,
		Body:     cleanBody(body),
		AuthorId: authorId,
	}
	if dbStructure.Chirps == nil {
		dbStructure.Chirps = map[int]Chirp{}
	}

	dbStructure.Chirps[chirpId] = newChirp
	err = db.writeDB(dbStructure)

	return newChirp, err
}

func (db *DB) DeleteChirp(chirpId int, userId int) error {
	err := db.ensureDB()
	if err != nil {
		return err
	}

	dbStructure, err := db.loadDB()

	if dbStructure.Chirps[chirpId].AuthorId != userId {
		return errors.New("Not authorized")
	}

	delete(dbStructure.Chirps, chirpId)
	db.writeDB(dbStructure)

	return nil
}

// // GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	err := db.ensureDB()
	if err != nil {
		return []Chirp{}, err
	}

	dbStructure, err := db.loadDB()
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, value := range dbStructure.Chirps {
		chirps = append(chirps, value)
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id > chirps[i].Id })
	return chirps, nil
}

func (db *DB) GetUserChirps(userId int, sortMethod string) ([]Chirp, error) {

	err := db.ensureDB()
	if err != nil {
		return []Chirp{}, err
	}

	dbStructure, err := db.loadDB()
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, value := range dbStructure.Chirps {
		if value.AuthorId == userId {
			chirps = append(chirps, value)
		}
	}

	if sortMethod == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id > chirps[i].Id })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id > chirps[i].Id })
	}
	return chirps, nil
}
func (db *DB) GetChirp(id int) (Chirp, error) {
	err := db.ensureDB()
	if err != nil {
		return Chirp{}, err
	}

	dbStructure, err := db.loadDB()
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirp, exists := dbStructure.Chirps[id]
	if exists {
		return chirp, nil
	} else {
		return Chirp{}, errors.New("does not exist")
	}
}

// // ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		empty_db, _ := json.Marshal(new(DBStructure))

		err := os.WriteFile(db.path, empty_db, 0666)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	err := db.ensureDB()
	if err != nil {
		return DBStructure{}, err
	}
	db.mux.Lock()
	defer db.mux.Unlock()

	db_data, _ := os.ReadFile(db.path)
	decoder := json.NewDecoder(bytes.NewReader(db_data))

	db_structure := DBStructure{}
	err = decoder.Decode(&db_structure)
	if err != nil {
		return db_structure, err
	}

	return db_structure, nil
}

// // writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

	err := db.ensureDB()
	if err != nil {
		return err
	}
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
