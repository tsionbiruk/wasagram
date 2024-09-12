package database

import "fmt"

//usercreation is in the tokens file.

func (db *wasabase) UpdateUserName(username string, newusername string) error {
	// Check if the new username already exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE username=?)", newusername).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if new username exists: %w", err)
	}
	if exists {
		return fmt.Errorf("new username %s already exists", newusername)
	}

	// Perform the update
	result, err := db.c.Exec("UPDATE Users SET username=? WHERE username=?", newusername, username)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
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
