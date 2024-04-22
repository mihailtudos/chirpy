package database

import (
	"errors"
	"time"
)

type Token string

// NotAuthorized error
var ErrNotAuthorized error = errors.New("not authorized")

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
