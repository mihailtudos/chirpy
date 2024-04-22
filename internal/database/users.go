package database

import (
	"github.com/mihailtudos/chirpy/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(user User) (User, error) {
	dbStr, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStr.Users) + 1
	user.ID = id
	dbStr.Users[id] = user
	if err := db.writeDB(dbStr); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStr, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users := dbStr.Users
	for _, u := range users {
		if u.Email == email {
			return u, nil
		}
	}

	return User{}, NotFound{}
}

func (db *DB) UpdateUser(userId int, email string, password string) (User, error) {
	dbStr, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStr.Users[userId]
	if !ok {
		return User{}, NotFound{}
	}

	if utils.IsPasswordValid(password) {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return User{}, err
		}
		user.Password = string(hashedPass)
	}

	if email != "" {
		user.Email = email
	}

	dbStr.Users[userId] = user

	return user, db.writeDB(dbStr)
}
