package database

import (
	"fmt"
)

func (db *wasabase) GetStream(username string) ([]photo, []string, error) {
	photos := []photo{}
	following := []string{}

	followRows, err := db.c.Query("SELECT username FROM Followes WHERE target_username = ?", username)
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
			err := photoRows.Scan(&p.photoId, &p.photo_png, &p.caption, &p.upload_time)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to scan photo: %w", err)
			}

			// Count the number of comments
			err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PhotoId=?", p.photoId).Scan(&p.comment_count)
			if err != nil {
				fmt.Println("Error executing comment count query:", err)
			}

			// Fetch all comments for the photo
			comments := []CommentData{}
			commentRows, err := db.c.Query("SELECT body, upload_time FROM Comments WHERE PhotoId = ?", p.photoId)
			if err != nil {
				return nil, nil, err
			}
			defer commentRows.Close()

			for commentRows.Next() {
				var comment CommentData
				err := commentRows.Scan(&comment.body, &comment.upload_time)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to scan comment: %w", err)
				}
				comments = append(comments, comment)
			}

			if err := commentRows.Err(); err != nil {
				return nil, nil, fmt.Errorf("error during comment row iteration: %w", err)
			}

			p.Comments = comments

			// Count the number of likes
			err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PhotoId=?", p.photoId).Scan(&p.like_count)
			if err != nil {
				fmt.Println("Error executing like count query:", err)
			}

			// Fetch all users who liked the photo
			likes := []string{}
			likeRows, err := db.c.Query("SELECT username FROM Likes WHERE PhotoId = ?", p.photoId)
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

			p.likes = likes

			// Append the photo to the photos slice
			photos = append(photos, p)
		}

		if err := photoRows.Err(); err != nil {
			return nil, nil, fmt.Errorf("error during photo row iteration: %w", err)
		}
	}

	return photos, following, nil
}
