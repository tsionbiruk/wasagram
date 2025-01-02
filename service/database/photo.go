package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) PhotoInsert(username string, photo []byte) error {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("INSERT INTO Photos (photo_id, user_id, photo_png,  upload_time) VALUES (NULL, ?, ?, ?)", user_id, photo, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) PhotoDelete(username string, photo_id int64) error {
	_, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return err
	}

	var exists bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Photos WHERE photo_id=?)", photo_id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check photo existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no photo found with PhotoId: %d", photo_id)
	}

	_, err = db.c.Exec("DELETE FROM Photos WHERE photo_id=?", photo_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) PhotoGet(photo_id int64) ([]byte, error) {
	var photo []byte
	err := db.c.QueryRow("SELECT photo_png FROM Photos WHERE photo_id=?", photo_id).Scan(&photo)
	if err != nil {
		return nil, err
	}
	return photo, nil
}
