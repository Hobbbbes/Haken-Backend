package handels

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/poodlenoodle42/Hacken-Backend/database"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	groups, err := database.GetGroupsForUser(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	vars := mux.Vars(r)
	groupIDstring := vars["groupID"]
	groupID, err := strconv.Atoi(groupIDstring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	tasks, err := database.GetTasksForGroup(token, groupID)
	if err != nil {
		if err.Error() == "User not allowed to view Group details" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		log.Println("GetTask: " + err.Error())
		return
	}
	json.NewEncoder(w).Encode(tasks)
}
