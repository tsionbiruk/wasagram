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

func (db *wasabase) UserProfile(username string, requester string) (*UserProfileInfo, error) {
	banned := []string{}

	rows, err := db.c.Query("SELECT target_username FROM Bans WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var users string
		if err := rows.Scan(&users); err != nil {
			return nil, fmt.Errorf("failed to scan banned: %w", err)
		}
		banned = append(banned, users)
	}
	fmt.Print(requester)

	for _, v := range banned {

		if v == requester {
			message := fmt.Sprintf("Requester '%s' is banned.", requester)
			return nil, fmt.Errorf(message)
		}
	}

	var userProfile UserProfileInfo
	userProfile.Username = username

	// Query the profile picture
	err = db.c.QueryRow("SELECT profil_pic FROM Users WHERE username = ?", username).Scan(&userProfile.ProfilPic)

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
	fmt.Print(userProfile.Follower_count)

	err = db.c.QueryRow("SELECT COUNT(*) FROM Followes WHERE username=?", username).Scan(&userProfile.Following_count)

	if err != nil {
		fmt.Println("Error executing query:", err)

	}
	fmt.Print(userProfile.Following_count)

	err = db.c.QueryRow("SELECT COUNT(*) FROM Bans WHERE username=?", username).Scan(&userProfile.Banned_count)

	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	fmt.Print(userProfile.Banned_count)

	//get photos
	photos := []photo{}
	rows, err = db.c.Query("SELECT photoId,photo_png,caption,upload_time FROM Photos WHERE username =? ORDER BY upload_time DESC", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p photo

		err := rows.Scan(&p.PhotoId, &p.Photo_png, &p.Caption, &p.Upload_time)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}

		err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PhotoId=?", p.PhotoId).Scan(&p.Comment_count)

		if err != nil {
			fmt.Println("Error executing query:", err)
		}

		Comments := []CommentData{}
		commentRows, err := db.c.Query("SELECT username, body, upload_time FROM Comments WHERE PhotoId=?", p.PhotoId)
		if err != nil {
			return nil, err
		}
		defer commentRows.Close()

		for commentRows.Next() {

			var comment CommentData
			err := commentRows.Scan(&comment.Author, &comment.Body, &comment.Upload_time)
			if err != nil {
				return nil, fmt.Errorf("failed to scan comment: %w", err)
			}
			Comments = append(Comments, comment)
		}
		if err := commentRows.Err(); err != nil {
			return nil, fmt.Errorf("error during comment row iteration: %w", err)
		}
		p.Comments = Comments
		fmt.Print(Comments)

		likes := []string{}

		err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PhotoId=?", p.PhotoId).Scan(&p.Like_count)

		if err != nil {
			fmt.Println("Error executing query:", err)
		}
		likeRows, err := db.c.Query("SELECT username FROM Likes WHERE photoId = ?", p.PhotoId)
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
		p.Likes = likes

		photos = append(photos, p)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	userProfile.Photo = photos

	userProfile.Photo_count = len(photos)

	return &userProfile, nil

}
