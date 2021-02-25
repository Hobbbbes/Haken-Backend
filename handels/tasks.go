package handels

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/poodlenoodle42/Hacken-Backend/database"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	groups, err := database.GetTasksForUser(token)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	tasks, err := database.GetTasksForUser(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("GetTask: " + err.Error())
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {

}
