package handels

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/poodlenoodle42/Hacken-Backend/container"
	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//Languages connects a language abbreviation to the language
var Languages map[string]datastructures.Language

func SubmitCode(w http.ResponseWriter, r *http.Request) {
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
	task, err := database.GetTask(taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	sub := datastructures.Submission{
		Author:  token,
		TaskID:  taskID,
		GroupID: task.GroupID,
	}
	sub, err = database.AddSubmission(sub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	langAb := r.FormValue("language")
	lang, ex := Languages[langAb]
	if !ex {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Language not known"))
		return
	}
	//inf, err := os.OpenFile(DataDir+fmt.Sprintf("/subtasks/%d_in", sub.ID), os.O_WRONLY|os.O_CREATE, 0666)
	path := DataDir + fmt.Sprintf("/submission/%d_in.%s", sub.ID, lang.Abbreviation)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	code := r.FormValue("code")
	f.Write([]byte(code))
	instance := container.GetInstance()
	defer container.ReturnInstance(instance)

	subtasks, err := database.GetSubtasksForTask(taskID, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	status, err := container.PrepareExecution(path, lang, instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if status.ExitCode != -1 {
		var preLaunchSubtask datastructures.Subtask
		for _, subT := range subtasks {
			if subT.Name == "PreLaunch" {
				preLaunchSubtask = subT
				break
			}
		}
		preLaunchRes := datastructures.Result{
			Sub:     &sub,
			Subt:    &preLaunchSubtask,
			Success: status.ExitCode,
		}
		err = database.AddResult(preLaunchRes)
		//TODO: Send Result via Websocket
	}
	for _, subtask := range subtasks {
		if subtask.Name == "PreLaunch" {
			continue
		}
		outPath := DataDir + fmt.Sprintf("/subtasks/%d_out", sub.ID)
		inPath := DataDir + fmt.Sprintf("/subtasks/%d_in", sub.ID)
		expectedOutBytes, err := ioutil.ReadFile(outPath)
		if err != nil {
			log.Println("SubmissionHandler: " + err.Error())
		}
		expectedOut := string(expectedOutBytes)
		inF, err := os.Open(inPath)
		defer inF.Close()
		s, err := container.Exec(instance, lang.LaunchTask, inF)
		if err != nil {
			log.Println("SubmissionHandler: " + err.Error())
		}
		Res := datastructures.Result{
			Sub:     &sub,
			Subt:    &subtask,
			Success: status.ExitCode,
		}
		if s.Output == expectedOut {
			Res.Success = -1
		}
		err = database.AddResult(Res)
		if err != nil {
			log.Println("SubmissionHandler: " + err.Error())
		}
		//TODO: Send result via websocket
	}
}
