package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/poodlenoodle42/Hacken-Backend/config"
	"github.com/poodlenoodle42/Hacken-Backend/database"
	"github.com/poodlenoodle42/Hacken-Backend/handels"
)

func main() {
	config := config.ReadConfig("config/config.yaml")
	handels.DataDir = config.DataDir
	f, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	database.InitDB(config.DBName, config.DBUser, config.DBPassword)
	defer database.CloseDB()

	r := mux.NewRouter().StrictSlash(true)
	//Use for unautherized route
	r.HandleFunc("/register", handels.AddUser).Methods("POST")
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(handels.AuthToken)
	s.HandleFunc("/groups", handels.GetGroups).Methods("GET")
	s.HandleFunc("/groups/{groupID}/rqtoken", handels.RequestGroupToken).Methods("GET")
	s.HandleFunc("/groups/{groupID}/tasks", handels.GetTasks).Methods("GET")
	s.HandleFunc("/groups/{groupID}/newTask", handels.NewTask).Methods("POST")
	s.HandleFunc("/groups/new", handels.NewGroup).Methods("POST")

	//Token as json
	s.HandleFunc("/groups/join", handels.JoinGroup).Methods("POST")

	s.HandleFunc("/tasks/{taskID}/subtasks", handels.GetSubtasks).Methods("GET")
	s.HandleFunc("/tasks/{taskID}", handels.GetTask).Methods("GET")
	s.HandleFunc("/tasks", handels.GetAllTasksForUser).Methods("GET")
	fmt.Println("Started serving")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
	//End

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
