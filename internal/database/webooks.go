package database

func (db *DB) UpgradedUserToRedChirpy(UserID int) error {
	dbStr, err := db.loadDB()
	if err != nil {
		return err
	}

	user := dbStr.Users[UserID]

	user.IsChirpyRed = true

	dbStr.Users[UserID] = user

	return db.writeDB(dbStr)
}
