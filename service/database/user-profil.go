package database

import (
	"database/sql"

	"fmt"
)

func (db *wasabase) UserProfile(username string) (*UserProfileInfo, error) {
	// get profil_pic

	row := db.c.QueryRow("SELECT username, profil_pic FROM Users WHERE username = '%s'", username)

	// Define variables to hold the query result
	var userProfile UserProfileInfo
	var profilPic []byte

	// Scan the result into variables
	err := row.Scan(userProfile.username, &profilPic)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with username: %s", username)
		}
		return nil, fmt.Errorf("failed to query user profile: %w", err)
	}

	// Assign the profile picture
	userProfile.profilPic = profilPic

	//get followers
	followers := []string{}
	rows, err := db.c.Query("SELECT target_username FROM Followes WHERE username = %s", username)
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

	//get following
	following := []string{}
	rows, err = db.c.Query("SELECT username FROM Followes WHERE target_username = %s", username)
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

	//get banned users

	banned := []string{}
	rows, err = db.c.Query("SELECT target_username FROM Bans WHERE username = %s", username)
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

	//get photos

	photos := []photo{}
	rows, err = db.c.Query("SELECT photoId,photo_png,caption,upload_time FROM Photos WHERE username = %s", username)
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

		comments := []CommentData{}
		commentRows, err := db.c.Query("SELECT body, upload_time FROM Comments WHERE photoId: %s", p.photoId)
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

		likes := []string{}
		likeRows, err := db.c.Query("SELECT username FROM Likes WHERE photoId = ?", p.photoId)
		if err != nil {
			return nil, err
		}

		defer likeRows.Close()
		for likeRows.Next() {
			var user string
			err := rows.Scan(&user)
			if err != nil {
				return nil, err
			}
			likes = append(likes, user)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		p.likes = likes

		photos = append(photos, p)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	userProfile.photo = photos

	return &UserProfileInfo{username, profilPic, followers, following, banned, photos}, nil

}
