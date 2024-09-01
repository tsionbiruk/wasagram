package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// uploadphoto, deletephoto, getphoto, likephoto, unlikephoto, comment, uncomment

// figure out how to get username of the logged in user.
func (db *wasabase) GetAuthorId(PhotoId int64) (string, error) {
	var author_name string
	err := db.c.QueryRow("SELECT username FROM Photos WHERE PhotoId=?", PhotoId).Scan(&author_name)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no row was found
			return "", fmt.Errorf("no author found for PhotoId %d", PhotoId)
		}
		// For other errors, return them
		return "", fmt.Errorf("error executing query: %v", err)
	}

	return author_name, nil
}

func (db *wasabase) Photolike(username string, PhotoId int64) error {
	err := db.c.QueryRow("SELECT * FROM Likes WHERE username=? AND PhotoId=?", username, PhotoId).Scan()
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.c.Exec("INSERT INTO Likes (username, PhotoId) VALUES (?, ?)", username, PhotoId)
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) Photounlike(username string, PhotoId int64) error {
	_, err := db.c.Exec("DELETE FROM Likes WHERE username=? AND PhotoId=?", username, PhotoId)
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) UploadPhoto(username string, caption string, photo []byte) error {

	_, err := db.c.Exec("INSERT INTO Photos (PhotoId, username , photo_png, caption, upload_time) VALUES (NULL,?, ?, ?, ?)", username, photo, caption, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) DeletePost(PhotoId int64) error {
	// Check if the photo exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Photos WHERE PhotoId=?)", PhotoId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check photo existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no photo found with PhotoId: %d", PhotoId)
	}

	// Delete the photo
	_, err = db.c.Exec("DELETE FROM Photos WHERE PhotoId=?", PhotoId)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	return nil
}

func (db *wasabase) PhotoGet(PhotoId int64) ([]byte, error) {
	var photo []byte
	err := db.c.QueryRow("SELECT photo_png FROM Photos WHERE PhotoId=?", PhotoId).Scan(&photo)
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (db *wasabase) comment(username string, PhotoId int64, text string) error {
	// Check if the photo exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Photos WHERE PhotoId=?)", PhotoId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check photo existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no photo found with PhotoId: %d", PhotoId)
	}

	// Insert the new comment
	_, err = db.c.Exec("INSERT INTO Comments (CommentId, username, PhotoId, body, upload_time) VALUES (NULL, ?, ?, ?, ?)",
		username, PhotoId, text, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

func (db *wasabase) uncomment(username string, PhotoId int64, CommentId int64) error {
	// Check if the comment exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comments WHERE CommentId=? AND PhotoId=? AND username=?)", CommentId, PhotoId, username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no comment found with comment_id: %d", CommentId)
	}

	// Delete the comment
	_, err = db.c.Exec("DELETE FROM Comments WHERE CommentId=? AND PhotoId=? AND username=?", CommentId, PhotoId, username)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}
