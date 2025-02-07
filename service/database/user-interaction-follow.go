package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) UserFollow(user_id int64, target_id int64) error {
	if user_id != target_id {
		var exists bool
		err := db.c.QueryRow("SELECT 1 FROM Follows WHERE user_id=? AND target_id=?", user_id, target_id).Scan(&exists)
		if !errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		var banned bool
		err = db.c.QueryRow("SELECT 1 FROM Bans WHERE user_id=? AND target_id=?", target_id, user_id).Scan(&banned)
		if !errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		_, err = db.c.Exec("INSERT INTO Follows (user_id, target_id) VALUES (?, ?)", user_id, target_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) UserUnfollow(user_id int64, target_id int64) error {
	if user_id != target_id {
		var exists bool
		err := db.c.QueryRow("SELECT 1 FROM Users WHERE user_id=?", target_id).Scan(&exists)
		if !exists {

			return nil
		} else if err != nil {

			return err
		}

		_, err = db.c.Exec("DELETE FROM Follows WHERE user_id=? AND target_id=?", user_id, target_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) UserGetFollowed(username string) ([]string, error) {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return nil, err
	}

	followed := make([]string, 0)
	rows, err := db.c.Query("SELECT user_name FROM Follows a INNER JOIN Users b ON a.target_id=b.user_id WHERE a.user_id=?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followed_name string
		err := rows.Scan(&followed_name)
		if err != nil {
			return nil, err
		}
		followed = append(followed, followed_name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followed, nil
}
