package handels

import (
	"net/http"
	"strings"

	"github.com/poodlenoodle42/Hacken-Backend/database"
)

var DataDir string

//AuthToken Authenticates a token
func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		token = strings.TrimSpace(token)
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		valid := database.AuthToken(token)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
