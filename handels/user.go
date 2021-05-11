package handels

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//AddUser checks if User exists and if not, adds user
func AddUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var user datastructures.UserLogin
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = database.AddUser(user)
	if err != nil {
		if err.Error() == "User already exists" {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
}
