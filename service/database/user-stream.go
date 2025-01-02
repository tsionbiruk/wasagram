package database

import (
	"strings"
)

func (db *appdbimpl) UserStream(username string) ([]StreamPost, error) {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return nil, err
	}

	following := make([]int64, 0)
	rows, err := db.c.Query("SELECT target_id FROM Follows WHERE user_id=?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var following_id int64
		err := rows.Scan(&following_id)
		if err != nil {
			return nil, err
		}
		following = append(following, following_id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	banned_by := make([]int64, 0)
	rows, err = db.c.Query("SELECT user_id FROM Bans WHERE target_id=?", user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var banner_id int64
		err := rows.Scan(&banner_id)
		if err != nil {
			return nil, err
		}
		banned_by = append(banned_by, banner_id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	photos := make([]StreamPost, 0)
	var followingString, bannedString string
	switch len(following) {
	case 0:
		followingString = ""
	case 1:
		followingString = "?"
	default:
		followingString = strings.Repeat("?,", len(following)-1) + "?"
	}
	switch len(banned_by) {
	case 0:
		bannedString = ""
	case 1:
		bannedString = "?"
	default:
		bannedString = strings.Repeat("?,", len(banned_by)-1) + "?"
	}
	values := make([]interface{}, len(following)+len(banned_by))
	for i, v := range following {
		values[i] = v
	}
	for i, v := range banned_by {
		values[len(following)+i] = v
	}
	rows, err = db.c.Query("SELECT photo_id, user_name, upload_time FROM Photos a INNER JOIN Users b ON a.user_id=b.user_id WHERE a.user_id IN ("+followingString+") AND a.user_id NOT IN ("+bannedString+") ORDER BY upload_time DESC", values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo_id int64
		var upload_time int64
		var author_username string

		err := rows.Scan(&photo_id, &author_username, &upload_time)
		if err != nil {
			return nil, err
		}

		comments := make([]CommentData, 0)
		rows, err := db.c.Query("SELECT comment_id, user_name, body, upload_time FROM Comments a INNER JOIN Users b ON a.user_id=b.user_id WHERE a.photo_id=? ORDER BY upload_time DESC", photo_id)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var comment_id int64
			var author_username string
			var content string
			var upload_time int64
			err := rows.Scan(&comment_id, &author_username, &content, &upload_time)
			if err != nil {
				return nil, err
			}
			comments = append(comments, CommentData{comment_id, author_username, upload_time, content})
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

		likes := make([]string, 0)
		rows, err = db.c.Query("SELECT user_name FROM Likes a INNER JOIN Users b ON a.user_id=b.user_id WHERE a.photo_id=?", photo_id)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var author_username string
			err := rows.Scan(&author_username)
			if err != nil {
				return nil, err
			}
			likes = append(likes, author_username)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

		photos = append(photos, StreamPost{photo_id, author_username, upload_time, likes, comments})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
