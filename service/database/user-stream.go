package database

import (
	"fmt"
)

func (db *wasabase) GetStream(username string, requester string) ([]photo, []string, error) {

	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE username=?", username).Scan(&exists)
	if !exists {

		return nil, nil, fmt.Errorf("user %s doesnt exist", username)
	} else if err != nil {

		return nil, nil, fmt.Errorf("error getting banned: %w", err)
	}

	banned := []string{}

	rows, err := db.c.Query("SELECT target_username FROM Bans WHERE username = ?", username)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var users string
		if err := rows.Scan(&users); err != nil {
			return nil, nil, fmt.Errorf("failed to scan following: %w", err)
		}
		banned = append(banned, users)
	}
	fmt.Print(requester)

	for _, v := range banned {

		if v == requester {
			message := fmt.Sprintf("Requester '%s' is banned.", requester)
			return nil, nil, fmt.Errorf(message)
		}
	}

	photos := []photo{}
	following := []string{}

	followRows, err := db.c.Query("SELECT target_username FROM Followes WHERE username = ?", username)
	if err != nil {
		return nil, nil, err
	}
	defer followRows.Close()

	for followRows.Next() {
		var follower string
		if err := followRows.Scan(&follower); err != nil {
			return nil, nil, fmt.Errorf("failed to scan following: %w", err)
		}
		following = append(following, follower)
	}

	if err := followRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error during row iteration: %w", err)
	}

	// For each user being followed, retrieve their photos
	for _, follower := range following {
		photoRows, err := db.c.Query("SELECT photoId, photo_png, caption, upload_time FROM Photos WHERE username = ? ORDER BY upload_time DESC", follower)
		if err != nil {
			return nil, nil, err
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var p photo
			err := photoRows.Scan(&p.PhotoId, &p.Photo_png, &p.Caption, &p.Upload_time)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to scan photo: %w", err)
			}

			// Count the number of comments
			err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PhotoId=?", p.PhotoId).Scan(&p.Comment_count)
			if err != nil {
				fmt.Println("Error executing comment count query:", err)
			}

			// Fetch all comments for the photo
			Comments := []CommentData{}
			commentRows, err := db.c.Query("SELECT body, upload_time FROM Comments WHERE PhotoId = ?", p.PhotoId)
			if err != nil {
				return nil, nil, err
			}
			defer commentRows.Close()

			for commentRows.Next() {
				var Comment CommentData
				err := commentRows.Scan(&Comment.Body, &Comment.Upload_time)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to scan comment: %w", err)
				}
				Comments = append(Comments, Comment)
			}

			if err := commentRows.Err(); err != nil {
				return nil, nil, fmt.Errorf("error during comment row iteration: %w", err)
			}

			p.Comments = Comments
			fmt.Print(Comments)
			// Count the number of likes
			err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PhotoId=?", p.PhotoId).Scan(&p.Like_count)
			if err != nil {
				fmt.Println("Error executing like count query:", err)
			}

			// Fetch all users who liked the photo
			likes := []string{}
			likeRows, err := db.c.Query("SELECT username FROM Likes WHERE PhotoId = ?", p.PhotoId)
			if err != nil {
				return nil, nil, err
			}
			defer likeRows.Close()

			for likeRows.Next() {
				var user string
				err := likeRows.Scan(&user)
				if err != nil {
					return nil, nil, err
				}
				likes = append(likes, user)
			}

			if err = likeRows.Err(); err != nil {
				return nil, nil, err
			}

			p.Likes = likes

			// Append the photo to the photos slice
			photos = append(photos, p)
		}

		if err := photoRows.Err(); err != nil {
			return nil, nil, fmt.Errorf("error during photo row iteration: %w", err)
		}
	}

	return photos, following, nil
}
