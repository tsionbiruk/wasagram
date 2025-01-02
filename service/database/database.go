/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateIfNoUser(username string) error
	UserRename(username string, newname string) error
	UserProfile(username string) (*UserProfileInfo, error)
	UserStream(username string) ([]StreamPost, error)
	GetAllUsers() ([]string, error)

	GetPhotoAuthorId(photo_id int64) (int64, error)
	GetUserIdFromUserName(username string) (int64, error)

	UserFollow(user_id int64, target_id int64) error
	UserUnfollow(user_id int64, target_id int64) error
	UserGetFollowed(username string) ([]string, error)

	UserBan(user_id int64, target_id int64) error
	UserUnban(user_id int64, target_id int64) error
	UserGetBanned(username string) ([]string, error)

	PhotoInsert(username string, photo []byte) error
	PhotoGet(photo_id int64) ([]byte, error)
	PhotoDelete(username string, photo_id int64) error

	PhotoLike(user_id int64, photo_id int64) error
	PhotoUnlike(user_id int64, photo_id int64) error

	PhotoComment(username string, photo_id int64, text string) error
	PhotoUncomment(username string, photo_id int64, comment_id int64) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	// Check if db is nil
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		usersTable := `CREATE TABLE Users (
			user_id INTEGER ,
			user_name TEXT UNIQUE,
			PRIMARY KEY(user_id)
		);`
		_, err = db.Exec(usersTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Follows';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		followsTable := `CREATE TABLE Follows (
			user_id INTEGER,
			target_id INTEGER,
			PRIMARY KEY(user_id, target_id),
			FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY(target_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
		_, err = db.Exec(followsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Bans';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		bansTable := `CREATE TABLE Bans (
			user_id INTEGER,
			target_id INTEGER,
			PRIMARY KEY(user_id, target_id),
			FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY(target_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
		_, err = db.Exec(bansTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		photosTable := `CREATE TABLE Photos (
			photo_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			photo_png BLOB,
			
			upload_time INTEGER,
			
			FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
		_, err = db.Exec(photosTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		likesTable := `CREATE TABLE Likes (
			user_id INTEGER,
			photo_id INTEGER,
			PRIMARY KEY(user_id, photo_id),
			FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY(photo_id) REFERENCES Photos(photo_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
		_, err = db.Exec(likesTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		commentsTable := `CREATE TABLE Comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			photo_id INTEGER,
			body TEXT,
			upload_time INTEGER,
			
			FOREIGN KEY(user_id) REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY(photo_id) REFERENCES Photos(photo_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`
		_, err = db.Exec(commentsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
