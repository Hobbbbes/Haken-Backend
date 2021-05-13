package handels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
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
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(groups)
}

func RequestGroupToken(w http.ResponseWriter, r *http.Request) {
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
	ex, err := database.DoesGroupExists(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !ex {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Group does not exist"))
		return
	}
	isAdmin, err := database.IsUserAdminOfGroup(token, groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User is not admin of group"))
		return
	}
	gToken := database.GenerateGroupToken(groupID)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte(fmt.Sprintf(`{"groupToken":"%s"}`, gToken)))

}

func JoinGroup(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var v interface{}
	err = json.Unmarshal(reqBody, &v)
	data := v.(map[string]interface{})
	gToken := fmt.Sprintf("%v", data["groupToken"])
	groupID := database.GetGroupIDFromToken(gToken)
	if groupID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Token does not exist"))
		return
	}
	err = database.AddUserToGroup(token, groupID)
	if err != nil {
		if err.Error() == "Group does not exists" || err.Error() == "User already in group" {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	group, err := database.GetGroup(groupID)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(err.Error()))
		return
	}
	tasks, err := database.GetTasksForGroup(group.ID)
	if err != nil {
		if err.Error() == "User not allowed to view Group details" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	res := datastructures.GroupWithTasks{
		group,
		tasks,
	}
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(res)
}

func NewGroup(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var group datastructures.Group
	err = json.Unmarshal(reqBody, &group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	group, err = database.AddNewGroup(token, group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(group)
}
