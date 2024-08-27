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
	GetName() (string, error)
	SetName(name string) error

	CreateNewUser(username string) error
	UpdateUserName(username string, newusername string) error
	GetUserIdFromUserName(username string) (int64, error)
	FollowUser(UserId int64, target_id int64) error
	UnFollowUser(UserId int64, target_id int64) error

	GetAllUsers() ([]string, error)

	BanUser(Userid int64, target_id int64) error
	UnBanUser(Userid int64, target_id int64) error
	GetBannedUsers(username string) ([]string, error)

	//other users
	GetUserFollowers(username string) ([]string, error)
	GetUserFollowing(username string) ([]string, error)
	GetProfile(username string) (*UserProfileInfo, error)
	GetStream(username string) ([]Posts, error)

	UploadPhoto(username string, photo []byte) error
	DeletePost(username string, postId int64) error //when you delete a photo you delete the comment and likes as well
	PhotoGet(PhotoId int64) ([]byte, error)
	Photolike(UserId int64, PhotoId int64) error
	Photounlike(UserId int64, PhotoId int64) error

	comment(username string, PostId int64, text string) error
	uncomment(username string, PostId int64, CommentId int64) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
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
			UserId INTEGER ,
			ProfilId INTEGER,
			username TEXT UNIQUE,
			PRIMARY KEY(UserId),
			FOREIGN KEY(ProfilId) REFERENCES Profils(ProfilId) ON DELETE CASCADE
		);`
		_, err = db.Exec(usersTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Followers';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		followersTable := `CREATE TABLE Followers(
			UserId INTEGER,
			target_id INTEGER,
			PRIMARY KEY(UserId, target_id),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(target_id) REFERENCES Users(UserId) ON DELETE CASCADE
		);`
		_, err = db.Exec(followersTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Following';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		followingTable := `CREATE TABLE Follows (
			UserId INTEGER,
			target_id INTEGER,
			PRIMARY KEY(UserId, target_id),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(target_id) REFERENCES Users(UserId) ON DELETE CASCADE
		);`
		_, err = db.Exec(followingTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Bans';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		bansTable := `CREATE TABLE Bans (
			UserId INTEGER,
			target_id INTEGER,
			PRIMARY KEY(UserId, target_id),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(target_id) REFERENCES Users(UserId) ON DELETE CASCADE
		);`
		_, err = db.Exec(bansTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		photosTable := `CREATE TABLE Photos (
			PhotoId INTEGER,
			UserId INTEGER,
			photo_png STRING,
			caption TEXT,
			upload_time DATE,
			PRIMARY KEY(PhotoId),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE
		);`
		_, err = db.Exec(photosTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		likesTable := `CREATE TABLE Likes (
			UserId INTEGER,
			PhotoTd INTEGER,
			PRIMARY KEY(UserId, PhotoId),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(PhotoId) REFERENCES Photos(PhotoId) ON DELETE CASCADE
		);`
		_, err = db.Exec(likesTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		commentsTable := `CREATE TABLE Comments (
			CommentId INTEGER,
			UserId INTEGER,
			PhotoId INTEGER,
			body TEXT,
			upload_time DATE,
			PRIMARY KEY(CommentId),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(PhotoId) REFERENCES Photos(PhotoId) ON DELETE CASCADE
		);`
		_, err = db.Exec(commentsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Posts';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		PostsTable := `CREATE TABLE Posts (
			PostId INTEGER,
			UserId INTEGER,
			PhotoId INTEGER,
			CommentId INTEGER,
			
			Nocomments INTEGER,
			Nolikes INTEGER,
			
			PRIMARY KEY(PostId, UserId, PhotoId, CommentId),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(PhotoId) REFERENCES Photos(PhotoId) ON DELETE CASCADE,
			FOREIGN KEY(CommentId) REFERENCES Comments(CommentId) ON DELETE CASCADE
		);`
		_, err = db.Exec(PostsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Profils';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		ProfilsTable := `CREATE TABLE Profils (
			ProfilId INTEGER,
			Profil_pic STRING,
			UserId INTEGER,
			Posts/stream ARRAY,
			PostId INTEGER,
			Followers INTEGER,
			Following INTEGER,
			BannedUsers INTEGER,
			PRIMARY KEY(UserId, ProfilId),
			FOREIGN KEY(UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
			FOREIGN KEY(PostId) REFERENCES Posts(PostId) ON DELETE CASCADE
			
		);`
		_, err = db.Exec(ProfilsTable)
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
