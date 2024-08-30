package database

import (
	"database/sql"

	"fmt"
)

// userinteractions: followuser, unfollowuser, banuser, unbanuser,
//getstream

func (db *wasabase) BanUsers(username string, target_username string) error {
	// Check if the ban already exists
	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Bans WHERE username=? AND target_username=?", username, target_username).Scan(&exists)
	if err == nil {
		// If no error and result is returned, it means the ban already exists
		return fmt.Errorf("user %s has already banned %s", username, target_username)
	} else if err != sql.ErrNoRows {
		// If any other error occurred
		return fmt.Errorf("failed to check ban status: %w", err)
	}

	// Proceed with banning the user
	_, err = db.c.Exec("INSERT INTO Bans (username, target_username) VALUES (?, ?)", username, target_username)
	if err != nil {
		return fmt.Errorf("failed to ban user %s: %w", target_username, err)
	}
	return nil
}

func (db *wasabase) UnBanUser(username string, target_username string) error {

	_, err := db.c.Exec("DELETE FROM Bans WHERE username=? AND target_username=?", username, target_username)
	if err != nil {
		return err
	}
	return nil
}
