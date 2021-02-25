package database

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

//Used when Password and Username are implemented
var mutex = &sync.Mutex{}
var recentlyUsedTokens map[string]bool

//AuthToken checks if a given token is in database
func AuthToken(token string) bool {
	a, ex := recentlyUsedTokens[token]
	if !ex {
		auth := authTokenDatabase(token)
		mutex.Lock()
		recentlyUsedTokens[token] = auth
		mutex.Unlock()
		time.AfterFunc(time.Hour, func() {
			mutex.Lock()
			delete(recentlyUsedTokens, token)
			mutex.Unlock()
		})
		return auth
	}
	return a
}
func authTokenDatabase(token string) bool {
	var exists bool
	err := db.QueryRow("SELECT exists (SELECT * FROM User WHERE Token = ?)", token).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false
	}
	return exists
}
