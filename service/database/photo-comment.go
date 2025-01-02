package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) PhotoComment(username string, photo_id int64, text string) error {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return err
	}

	// Check if the photo exists
	var exists bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Photos WHERE photo_id=?)", photo_id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check photo existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no photo found with PhotoId: %d", photo_id)
	}

	// Insert the new comment

	_, err = db.c.Exec("INSERT INTO Comments (comment_id, user_id, photo_id, body, upload_time) VALUES (NULL, ?, ?, ?, ?)", user_id, photo_id, text, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) PhotoUncomment(username string, photo_id int64, comment_id int64) error {
	_, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return err
	}

	var exists bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comments WHERE comment_id=?)", comment_id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no comment found with comment_id: %d", comment_id)
	}

	_, err = db.c.Exec("DELETE FROM Comments WHERE comment_id=?", comment_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) PhotoLike(user_id int64, photo_id int64) error {
	err := db.c.QueryRow("SELECT * FROM Likes WHERE user_id=? AND photo_id=?", user_id, photo_id).Scan()
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.c.Exec("INSERT INTO Likes (user_id, photo_id) VALUES (?, ?)", user_id, photo_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) PhotoUnlike(user_id int64, photo_id int64) error {
	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE user_id=?", user_id).Scan(&exists)
	if !exists {

		return fmt.Errorf("user %d doesnt exist", user_id)
	} else if err != nil {

		return fmt.Errorf("error banning follower: %w", err)
	}

	_, err = db.c.Exec("DELETE FROM Likes WHERE user_id=? AND photo_id=?", user_id, photo_id)
	if err != nil {
		return err
	}

	return nil
}
