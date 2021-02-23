package database

import "time"

//Used when Password and Username are implemented

var recentlyUsedTokens map[string]bool

//AuthToken checks if a given token is in database
func AuthToken(token string) bool {
	a, ex := recentlyUsedTokens[token]
	if !ex {
		auth := authTokenDatabase(token)
		recentlyUsedTokens[token] = auth
		time.AfterFunc(time.Hour, func() {
			delete(recentlyUsedTokens, token)
		})
		return auth
	}
	return a
}
func authTokenDatabase(token string) bool {
	return true
}
