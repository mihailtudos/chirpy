package database

import "fmt"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	if id < 1 {
		return Chirp{}, fmt.Errorf("the id must be a positive number")
	}

	dbStr, err := db.loadDB()
	if err != nil {
		fmt.Println(err.Error())
		return Chirp{}, fmt.Errorf("failed to load the DB")
	}

	chirp, ok := dbStr.Chirps[id]
	if !ok {
		return Chirp{}, NotFound{}
	}

	return chirp, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		fmt.Println("failed to load the db chirp records")
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}

	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		fmt.Println("could not write chirps on the db")
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
