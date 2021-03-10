package handels

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

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

	ex, err := database.DoesGroupExists(groupID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if !ex {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Group does not exist"))
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
		return
	}
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetSubtasks(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	vars := mux.Vars(r)
	taskIDstring := vars["taskID"]
	taskID, err := strconv.Atoi(taskIDstring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	allowed, err := database.IsUserAllowedToAccessTask(token, taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !allowed {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User not allowed to access task"))
		return
	}
	subtasks, err := database.GetSubtasksForTask(taskID, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(subtasks)
}

func GetTaskPDF(w http.ResponseWriter, r *http.Request) {
	//	token := r.Header.Get("token")
	//	token = strings.TrimSpace(token)
	token := r.URL.Query().Get("token")
	if token == "" || !database.AuthToken(token) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	taskIDstring := vars["taskID"]
	taskID, err := strconv.Atoi(taskIDstring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	allowed, err := database.IsUserAllowedToAccessTask(token, taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !allowed {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User is not allowed to view task"))
		return
	}
	path := DataDir + fmt.Sprintf("/tasks/%d.pdf", taskID)
	if _, err := os.Stat(path); err == nil || os.IsExist(err) {
		http.ServeFile(w, r, path)
		return
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("File does not exit"))
	}

}

func NewTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
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

	file, _, err := r.FormFile("pdf")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var taskInfo datastructures.Task
	taskInfo.Description = r.FormValue("description")
	taskInfo.Name = r.FormValue("name")

	admin, err := database.IsUserAdminOfGroup(token, groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !admin {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User is not admin of Group"))
		return
	}
	taskInfo.GroupID = groupID
	taskInfo.Author = token
	task, err := database.AddTask(taskInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	PreLaunchSubtask := datastructures.Subtask{
		Points: 0,
		Name:   "PreLaunch",
		TaskID: task.ID,
	}
	PreLaunchSubtask, err = database.AddSubtask(PreLaunchSubtask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	f, err := os.OpenFile(DataDir+fmt.Sprintf("/tasks/%d.pdf", task.ID), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.Header().Set("Content-Type", "text/json")
	json.NewEncoder(w).Encode(task)
}

func GetAllTasksForUser(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	groups, err := database.GetGroupsForUser(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	groupIDs := make([]interface{}, len(groups))
	for index, v := range groups {
		groupIDs[index] = v.ID
	}
	tasks, err := database.GetTasksForGroups(groupIDs)
	w.Header().Set("Content-Type", "text/json")
	groupsWithTasks := make([]datastructures.GroupWithTasks, len(groups))
	for i, group := range groups {
		for _, task := range tasks {
			if task.GroupID == group.ID {
				groupsWithTasks[i] = datastructures.GroupWithTasks{
					group,
					append(groupsWithTasks[i].Tasks, task),
				}
			}
			groupsWithTasks[i] = datastructures.GroupWithTasks{
				group,
				groupsWithTasks[i].Tasks,
			}
			if len(groupsWithTasks[i].Tasks) == 0 {
				groupsWithTasks[i].Tasks = make([]datastructures.Task, 0)
			}
		}
	}
	json.NewEncoder(w).Encode(groupsWithTasks)

}
func NewSubtask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)

	vars := mux.Vars(r)
	taskIDstring := vars["taskID"]
	taskID, err := strconv.Atoi(taskIDstring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	all, err := database.IsUserAllowedToAccessTask(token, taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !all {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User not allowed to access task"))
		return
	}
	isAuthor, err := database.IsUserAuthorOfTask(token, taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !isAuthor {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("User not author of task"))
		return
	}
	subtaskInfoJSON := r.FormValue("info")
	if subtaskInfoJSON == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No task info"))
		return
	}
	var subtaskInfo datastructures.Subtask
	err = json.Unmarshal([]byte(subtaskInfoJSON), &subtaskInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	infile, _, err := r.FormFile("in")
	defer infile.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	outfile, _, err := r.FormFile("out")
	defer outfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	subtaskInfo.TaskID = taskID
	sub, err := database.AddSubtask(subtaskInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	outf, err := os.OpenFile(DataDir+fmt.Sprintf("/subtasks/%d_out", sub.ID), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer outf.Close()
	io.Copy(outf, outfile)
	inf, err := os.OpenFile(DataDir+fmt.Sprintf("/subtasks/%d_in", sub.ID), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer outf.Close()
	io.Copy(inf, infile)
	w.WriteHeader(http.StatusOK)
}
