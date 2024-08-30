package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// all user related functions: createnewuser, updateusername,
//getallusers,

func (db *wasabase) CreateNewUser(username string) error {
	var Profil_pic []byte
	err := db.c.QueryRow("SELECT * FROM Users WHERE username='%s'", username).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.c.Exec("INSERT INTO Users VALUES ('%s',%w)", username, Profil_pic)
		if err != nil {
			return fmt.Errorf("failed to insert '%s' into the Users table: %w", username, err)
		}
	}
	return nil
}

func (db *wasabase) UpdateUserName(username string, newusername string) error {
	_, err := db.c.Exec("UPDATE Users SET username='%s' WHERE username='%s'", newusername, username)
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) GetAllUsers() ([]string, error) {
	// Execute the query
	rows, err := db.c.Query("SELECT username FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string

	// Iterate over the rows
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}

	// Check for any error during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
