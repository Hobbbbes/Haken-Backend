package handels

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

var DataDir string

var CookieToToken map[string]string

//AuthToken Authenticates a login cookie
func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Tempory_Login_Token")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cookie := c.Value
		cookie = strings.TrimSpace(cookie)
		if cookie == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		token, valid := CookieToToken[cookie]
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		r.Header.Add("token", token)
		next.ServeHTTP(w, r)
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var userLogin datastructures.UserLogin
	err = json.Unmarshal(reqBody, &userLogin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := database.AuthUser(userLogin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cookieValue := database.RandomString(50)
	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{Name: "Tempory_Login_Token", Value: cookieValue, Expires: expiration, Path: "/"}
	CookieToToken[cookieValue] = user.Token
	http.SetCookie(w, &cookie)

}
