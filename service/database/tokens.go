package database

import (
	"database/sql"
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

	//if user exists we skip the user incertion and jump here.
	token := db.tokenGen.GenerateUniqueToken(username)
	currentTime := time.Now().Format(time.RFC3339)

	_, err = db.c.Exec(`REPLACE INTO Tokens(username, token, time) VALUES (?, ?, ?)`, username, token, currentTime)
	if err != nil {
		log.Fatal(err)
	}

	return token, nil
}

func (db *wasabase) Gettoken(username string) (int64, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Tokens WHERE username=?)", username).Scan(&exists)
	if err != nil {

		return 0, fmt.Errorf("failed to check user existence in Tokens: %w", err)
	}

	var token int64
	err = db.c.QueryRow("SELECT Token FROM Tokens WHERE username=?", username).Scan(&token)
	if err != nil {

		return 0, fmt.Errorf("failed to check user Tokens: %w", err)
	}

	return token, nil
}

func (db *wasabase) Gettokentime(username string, token int64) (time.Time, error) {
	var tokentime time.Time
	err := db.c.QueryRow("SELECT time FROM Tokens WHERE username=? AND token=?", username, token).Scan(&tokentime)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, fmt.Errorf("no token found for user '%s': %w", username, err)
		}
		return time.Time{}, fmt.Errorf("failed to query token time: %w", err)
	}
	return tokentime, nil
}

func (db *wasabase) Istokenexpired(tokenTime time.Time) bool {
	// Get the current time
	currentTime := time.Now()

	// Define the expiration duration (1 hour)
	expirationDuration := time.Hour

	// Check if the current time is after the token time + expiration duration
	return currentTime.Sub(tokenTime) > expirationDuration
}
