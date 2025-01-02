package database

func (db *appdbimpl) UserProfile(username string) (*UserProfileInfo, error) {
	user_id, err := db.GetUserIdFromUserName(username)
	if err != nil {
		return nil, err
	}

	// get the followers list
	followers := make([]string, 0)
	rows, err := db.c.Query("SELECT user_name FROM Follows a INNER JOIN Users b ON a.user_id=b.user_id WHERE a.target_id=?", user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var followers_name string
		err := rows.Scan(&followers_name)
		if err != nil {
			return nil, err
		}
		followers = append(followers, followers_name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// get the following list
	following := make([]string, 0)
	rows, err = db.c.Query("SELECT user_name FROM Follows a INNER JOIN Users b ON a.target_id=b.user_id WHERE a.user_id=?", user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var following_name string
		err := rows.Scan(&following_name)
		if err != nil {
			return nil, err
		}
		following = append(following, following_name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// get the banned list
	banned := make([]string, 0)
	rows, err = db.c.Query("SELECT user_name FROM Bans a INNER JOIN Users b ON a.target_id=b.user_id WHERE a.user_id=?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banned_name string
		err := rows.Scan(&banned_name)
		if err != nil {
			return nil, err
		}
		banned = append(banned, banned_name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// get the photos list
	photos := make([]StreamPost, 0)
	rows, err = db.c.Query("SELECT photo_id, upload_time FROM Photos WHERE user_id=? ORDER BY upload_time DESC", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo_id int64
		var upload_time int64
		err := rows.Scan(&photo_id, &upload_time)
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
			var author_name string
			var body string
			var upload_time int64
			err := rows.Scan(&comment_id, &author_name, &body, &upload_time)
			if err != nil {
				return nil, err
			}
			comments = append(comments, CommentData{comment_id, author_name, upload_time, body})
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

		likes := make([]string, 0)
		rows, err = db.c.Query("SELECT user_name FROM Likes a INNER JOIN Users b ON a.user_id=b.user_id WHERE a.photo_id=?", photo_id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var author_name string
			err := rows.Scan(&author_name)
			if err != nil {
				return nil, err
			}
			likes = append(likes, author_name)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

		photos = append(photos, StreamPost{photo_id, username, upload_time, likes, comments})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &UserProfileInfo{followers, following, banned, photos}, nil
}
