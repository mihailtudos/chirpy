package database

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
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
