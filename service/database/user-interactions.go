package database

import (
	"database/sql"

	"fmt"
)

// userinteractions: followuser, unfollowuser, banuser, unbanuser,

func (db *wasabase) BanUsers(username string, target_username string) error {

	var exists bool
	err := db.c.QueryRow("SELECT username FROM Users WHERE username=? ", target_username).Scan(&exists)
	if err == nil {

		return fmt.Errorf("%s doesnt exist", target_username)
	} else if err != sql.ErrNoRows {

		return fmt.Errorf("failed to check ban status: %w", err)
	}

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

func (db *wasabase) FollowUser(username string, target_username string) error {

	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Followes WHERE username=? AND target_username=?", username, target_username).Scan(&exists)
	if err == nil {

		return fmt.Errorf("user %s already follows %s", username, target_username)
	} else if err != sql.ErrNoRows {

		return fmt.Errorf("failed to check following status: %w", err)
	}

	_, err = db.c.Exec("INSERT INTO Followes (username, target_username) VALUES (?, ?)", username, target_username)
	if err != nil {
		return fmt.Errorf("failed to follow user %s: %w", target_username, err)
	}
	return nil
}
func (db *wasabase) UnFollowUser(username string, target_username string) error {

	_, err := db.c.Exec("DELETE FROM Followes WHERE username=? AND target_username=?", username, target_username)
	if err != nil {
		return err
	}
	return nil
}
