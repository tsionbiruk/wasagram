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

func (db *wasabase) GetAuthorliker(PhotoId int64) ([]string, error) {
	var author_names []string

	// Execute the query expecting multiple rows
	rows, err := db.c.Query("SELECT username FROM Likes WHERE PhotoId=?", PhotoId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows are found
			return nil, fmt.Errorf("no authors found for PhotoId %d", PhotoId)
		}
		// For other errors, return them
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Iterate over the result set and append to the author_names slice
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		author_names = append(author_names, username)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return author_names, nil

}

func (db *wasabase) GetAuthorcommenter(PhotoId int64) ([]string, error) {
	var author_names []string

	// Execute the query expecting multiple rows
	rows, err := db.c.Query("SELECT username FROM Comments WHERE PhotoId=?", PhotoId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows are found
			return nil, fmt.Errorf("no authors found for PhotoId %d", PhotoId)
		}
		// For other errors, return them
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Iterate over the result set and append to the author_names slice
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		author_names = append(author_names, username)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return author_names, nil

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
	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE username=?", username).Scan(&exists)
	if !exists {

		return fmt.Errorf("user %s doesnt exist", username)
	} else if err != nil {

		return fmt.Errorf("error banning follower: %w", err)
	}

	_, err = db.c.Exec("DELETE FROM Likes WHERE username=? AND PhotoId=?", username, PhotoId)
	if err != nil {
		return err
	}

	return nil
}

func (db *wasabase) Getlikes(PhotoId int64) ([]string, error) {
	likes := []string{}
	likeRows, err := db.c.Query("SELECT username FROM Likes WHERE photoId = ?", PhotoId)
	if err != nil {
		return nil, err
	}

	defer likeRows.Close()
	for likeRows.Next() {
		var user string
		err := likeRows.Scan(&user)
		if err != nil {
			return nil, err
		}
		likes = append(likes, user)

	}
	if err = likeRows.Err(); err != nil {
		return nil, err
	}
	return likes, nil
}

func (db *wasabase) UploadPhoto(username string, caption string, photo []byte) error {

	_, err := db.c.Exec("INSERT INTO Photos (PhotoId, username , photo_png, caption, upload_time) VALUES (NULL,?, ?, ?, ?)", username, photo, caption, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) DeletePost(PhotoId int64) (string, error) {
	// Check if the photo exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Photos WHERE PhotoId=?)", PhotoId).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to check photo existence: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("no photo found with PhotoId: %d", PhotoId)
	}

	// Delete the photo
	_, err = db.c.Exec("DELETE FROM Photos WHERE PhotoId=?", PhotoId)
	if err != nil {
		return "", fmt.Errorf("failed to delete photo: %w", err)
	}

	return "photo deleted succesfully!", nil
}

func (db *wasabase) Comment(username string, PhotoId int64, text string) error {
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

func (db *wasabase) Uncomment(CommentId int64) error {
	// Check if the comment exists
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comments WHERE CommentId=?)", CommentId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no comment found with comment_id: %d", CommentId)
	}

	// Delete the comment
	_, err = db.c.Exec("DELETE FROM Comments WHERE CommentId=?", CommentId)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (db *wasabase) Getcomment(PhotoId int64) ([]string, error) {
	comments := []string{}
	commentRows, err := db.c.Query("SELECT body FROM Comments WHERE photoId = ?", PhotoId)
	if err != nil {
		return nil, err
	}

	defer commentRows.Close()
	for commentRows.Next() {
		var user string
		err := commentRows.Scan(&user)
		if err != nil {
			return nil, err
		}
		comments = append(comments, user)

	}
	if err = commentRows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
