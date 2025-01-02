package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) UserRename(username string, newname string) error {
	result, err := db.c.Exec("UPDATE Users SET user_name=? WHERE user_name=?", newname, username)
	if err != nil {
		return err
	}
	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated, username %s not found", username)
	}
	return nil
}

func (db *appdbimpl) CreateIfNoUser(username string) error {
	err := db.c.QueryRow("SELECT * FROM Users WHERE user_name=?", username).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.c.Exec("INSERT INTO Users VALUES (NULL, ?)", username)
		if err != nil {
			return fmt.Errorf("failed to insert '%s' into the Users table: %w", username, err)
		}
	}
	return nil
}

func (db *appdbimpl) GetAllUsers() ([]string, error) {
	users := make([]string, 0)
	rows, err := db.c.Query("SELECT user_name FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		users = append(users, username)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *appdbimpl) GetUserIdFromUserName(username string) (int64, error) {
	var user_id int64
	err := db.c.QueryRow("SELECT user_id from Users WHERE user_name=?", username).Scan(&user_id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return user_id, nil
}

func (db *appdbimpl) GetPhotoAuthorId(photo_id int64) (int64, error) {
	var owner_id int64
	err := db.c.QueryRow("SELECT user_id FROM Photos WHERE photo_id=?", photo_id).Scan(&owner_id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return owner_id, nil
}
