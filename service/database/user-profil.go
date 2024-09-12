package database

import (
	"database/sql"

	"fmt"
)

func (db *wasabase) Getbannedusers(username string) ([]string, error) {

	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE username=?", username).Scan(&exists)
	if !exists {

		return nil, fmt.Errorf("user %s doesnt exist", username)
	} else if err != nil {

		return nil, fmt.Errorf("error getting banned: %w", err)
	}

	banned := []string{}

	rows, err := db.c.Query("SELECT target_username FROM Bans WHERE username =?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ban string
		if err := rows.Scan(&ban); err != nil {
			return nil, fmt.Errorf("failed to scan following: %w", err)
		}
		banned = append(banned, ban)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	return banned, nil
}

func (db *wasabase) Getfollowers(username string) ([]string, error) {

	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE username=?", username).Scan(&exists)
	if !exists {

		return nil, fmt.Errorf("user %s doesnt exist", username)
	} else if err != nil {

		return nil, fmt.Errorf("error getting follower: %w", err)
	}

	followers := []string{}
	rows, err := db.c.Query("SELECT username FROM Followes WHERE target_username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	return followers, nil

}

func (db *wasabase) Getfollowing(username string) ([]string, error) {

	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Users WHERE username=?", username).Scan(&exists)
	if !exists {

		return nil, fmt.Errorf("user %s doesnt exist", username)
	} else if err != nil {

		return nil, fmt.Errorf("error getting following: %w", err)
	}
	following := []string{}
	rows, err := db.c.Query("SELECT target_username FROM Followes WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			return nil, fmt.Errorf("failed to scan following: %w", err)
		}
		following = append(following, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	return following, nil
}

func (db *wasabase) UserProfile(username string) (*UserProfileInfo, error) {
	var userProfile UserProfileInfo
	userProfile.Username = username

	// Query the profile picture
	err := db.c.QueryRow("SELECT profil_pic FROM Users WHERE username = ?", username).Scan(&userProfile.ProfilPic)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with username: %s", username)
		}
		return nil, fmt.Errorf("failed to query user profile: %w", err)
	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM Followes WHERE target_username=?", username).Scan(&userProfile.Follower_count)

	if err != nil {
		fmt.Println("Error executing query:", err)

	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM Followes WHERE username=?", username).Scan(&userProfile.Following_count)

	if err != nil {
		fmt.Println("Error executing query:", err)

	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM Bans WHERE username=?", username).Scan(&userProfile.Banned_count)

	if err != nil {
		fmt.Println("Error executing query:", err)
	}

	//get photos
	photos := []photo{}
	rows, err := db.c.Query("SELECT photoId,photo_png,caption,upload_time FROM Photos WHERE username =? ORDER BY upload_time DESC", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p photo

		err := rows.Scan(&p.photoId, &p.photo_png, &p.caption, &p.upload_time)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}

		err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PhotoId=?", p.photoId).Scan(&p.comment_count)

		if err != nil {
			fmt.Println("Error executing query:", err)
		}

		comments := []CommentData{}
		commentRows, err := db.c.Query("SELECT body, upload_time FROM Comments WHERE photoId=?", p.photoId)
		if err != nil {
			return nil, err
		}
		defer commentRows.Close()

		for commentRows.Next() {

			var comment CommentData
			err := commentRows.Scan(&comment.body, &comment.upload_time)
			if err != nil {
				return nil, fmt.Errorf("failed to scan comment: %w", err)
			}
			comments = append(comments, comment)
		}
		if err := commentRows.Err(); err != nil {
			return nil, fmt.Errorf("error during comment row iteration: %w", err)
		}
		p.Comments = comments

		fmt.Print(comments)
		fmt.Print(p.comment_count)

		likes := []string{}

		err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PhotoId=?", p.photoId).Scan(&p.like_count)

		if err != nil {
			fmt.Println("Error executing query:", err)
		}
		likeRows, err := db.c.Query("SELECT username FROM Likes WHERE photoId = ?", p.photoId)
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

		if err = rows.Err(); err != nil {
			return nil, err
		}
		p.likes = likes
		fmt.Print(likes)
		fmt.Print(p.like_count)

		photos = append(photos, p)
		fmt.Print(photos)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	userProfile.Photo = photos
	userProfile.Photo_count = len(photos)

	return &userProfile, nil

}
