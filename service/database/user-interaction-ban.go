package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) UserBan(user_id int64, target_id int64) error {
	if user_id != target_id {
		var exists bool
		err := db.c.QueryRow("SELECT 1 FROM Bans WHERE user_id=? AND target_id=?", user_id, target_id).Scan(&exists)
		if !errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if !exists {

			return fmt.Errorf("user %d doesnt exist", target_id)
		}
		_, err = db.c.Exec("INSERT INTO Bans (user_id, target_id) VALUES (?, ?)", user_id, target_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) UserUnban(user_id int64, target_id int64) error {

	_, err := db.c.Exec("DELETE FROM Bans WHERE user_id=? AND target_id=?", user_id, target_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UserGetBanned(username string) ([]string, error) {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return nil, err
	}

	banned := make([]string, 0)
	rows, err := db.c.Query("SELECT user_name FROM Bans a INNER JOIN Users b ON a.target_id=b.user_id WHERE a.user_id=?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banned_name string
		err := rows.Scan(&banned_name)
		if err != nil {
			return nil, err
		}
		banned = append(banned, banned_name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banned, nil
}
