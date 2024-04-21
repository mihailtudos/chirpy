package database

import "time"

type Token string

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	dbStr, err := db.loadDB()
	if err != nil {
		return true, err
	}

	_, ok := dbStr.RevokedTokens[Token(token)]
	if !ok {
		return false, nil
	}

	return true, nil
}

func (db *DB) RevokeToken(token Token) error {
	dbStr, err := db.loadDB()
	if err != nil {
		return err
	}

	dbStr.RevokedTokens[token] = time.Now()

	return db.writeDB(dbStr)
}
