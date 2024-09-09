package database

//usercreation is in the tokens file.

func (db *wasabase) UpdateUserName(username string, newusername string) error {
	_, err := db.c.Exec("UPDATE Users SET username=? WHERE username=?", newusername, username)
	if err != nil {
		return err
	}
	return nil
}

func (db *wasabase) GetAllUsers() ([]string, error) {
	// Execute the query
	rows, err := db.c.Query("SELECT username FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string

	// Iterate over the rows
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}

	// Check for any error during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
