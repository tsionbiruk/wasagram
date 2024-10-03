package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"sync"
	"time"
)

// create structture to store tokens counter is a global
type TokenGenerator struct {
	mu      sync.Mutex
	counter int
}

// NewTokenGenerator creates and returns a new TokenGenerator creates new insatnce of
// token generator
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{
		counter: 0, // Start the global token counter from 0
	}
}

// GenerateUniqueToken generates and returns a globally unique token.
func (tg *TokenGenerator) GenerateUniqueToken(username string) int {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	// Increment the global counter to generate a unique token
	tg.counter++
	return tg.counter
}

// NewWasabase initializes the wasabase struct with a database connection and a TokenGenerator.
func NewWasabase(db *sql.DB) *wasabase {
	return &wasabase{
		c:        db,
		tokenGen: NewTokenGenerator(), // Initialize TokenGenerator here
	}
}

func (db *wasabase) Creatuser_Getuserfromtoken(username string) (int, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE username=?)", username).Scan(&exists)
	if err != nil {

		return 0, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		var Profil_pic []byte
		_, err = db.c.Exec("INSERT INTO Users VALUES (?,?)", username, Profil_pic)
		if err != nil {
			return 0, fmt.Errorf("failed to insert '%s' into the Users table: %w", username, err)

		}
	}

	// if user exists we skip the user incertion and jump here.
	token := db.tokenGen.GenerateUniqueToken(username)
	currentTime := time.Now().Format("15:04:05")

	_, err = db.c.Exec(`DELETE FROM Tokens WHERE username=?`, username)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.c.Exec(`REPLACE INTO Tokens(username, token, time) VALUES (?, ?, ?)`, username, token, currentTime)
	if err != nil {
		log.Fatal(err)
	}

	return token, nil
}

func (db *wasabase) Gettoken(username string) (int64, error) {
	var token int64
	err := db.c.QueryRow("SELECT Token FROM Tokens WHERE username=?", username).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no rows are found
			// For example, you might want to return a specific value or nil
			return 0, nil // Return 0 or another appropriate value indicating no result
		}
		// Handle other types of errors
		return 0, fmt.Errorf("failed to check user Tokens: %w", err)
	}
	return token, nil
}

func (db *wasabase) Gettokentime(username string, token int64) (time.Time, error) {
	var timeString string
	err := db.c.QueryRow("SELECT time FROM Tokens WHERE username=? AND Token=?", username, token).Scan(&timeString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil // No result found
		}
		return time.Time{}, fmt.Errorf("failed to query time: %w", err)
	}

	// Parse the time string to time.Time
	timeLayout := "15:04:05" // Adjust this layout to match your date format
	parsedTime, err := time.Parse(timeLayout, timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return parsedTime, nil
}

func (db *wasabase) Istokenexpired(tokenTime time.Time) bool {
	// Get the current time
	timeLayout := "15:04:05"
	currentTime := time.Now().Format("15:04:05")

	// Adjust this layout to match your date format
	parsedTime, err := time.Parse(timeLayout, currentTime)
	if err != nil {
		return false
	}

	// Define the expiration duration (1 hour)
	expirationDuration := time.Hour

	// Check if the current time is after the token time + expiration duration
	return parsedTime.Sub(tokenTime) > expirationDuration
}
