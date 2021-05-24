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
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
	"github.com/poodlenoodle42/Hacken-Backend/handels"
	"github.com/rs/cors"
)

func main() {
	config := config.ReadConfig("config/config.yaml")
	handels.DataDir = config.DataDir
	handels.Languages = make(map[string]datastructures.Language)
	handels.CookieToToken = make(map[string]string)
	for _, lang := range config.ContainerConfig.Languages {
		handels.Languages[lang.Abbreviation] = lang
	}
	f, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	database.InitDB(config.DBName, config.DBUser, config.DBPassword)
	defer database.CloseDB()

	//container.InitInstances(config.ContainerConfig)
	//defer container.StopAndDeleteInstances()
	r := mux.NewRouter().StrictSlash(true)
	//Use for unautherized route
	r.HandleFunc("/register", handels.AddUser).Methods("POST")
	r.HandleFunc("/login", handels.Login).Methods("POST")
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(handels.AuthToken)
	s.HandleFunc("/groups", handels.GetGroups).Methods("GET")
	s.HandleFunc("/groups/{groupID}/rqtoken", handels.RequestGroupToken).Methods("GET")
	//s.HandleFunc("/groups/{groupID}/tasks", handels.GetTasks).Methods("GET")
	s.HandleFunc("/groups/{groupID}/newTask", handels.NewTask).Methods("POST")
	s.HandleFunc("/groups/new", handels.NewGroup).Methods("POST")

	//Token as json
	s.HandleFunc("/groups/join", handels.JoinGroup).Methods("POST")

	//s.HandleFunc("/tasks/{taskID}/subtasks", handels.GetSubtasks).Methods("GET")
	r.HandleFunc("/auth/tasks/{taskID}/pdf", handels.GetTaskPDF).Methods("GET")
	s.HandleFunc("/tasks/{taskID}", handels.GetSubtasks).Methods("GET")
	s.HandleFunc("/tasks", handels.GetAllTasksForUser).Methods("GET")

	s.HandleFunc("/tasks/{taskID}/newSubtask", handels.NewSubtask).Methods("POST")
	s.HandleFunc("/tasks/{taskID}/delete", handels.DeleteTask).Methods("DELETE")
	s.HandleFunc("/tasks/{taskID}/submit", handels.SubmitCode).Methods("POST")
	s.HandleFunc("/subtasks/{subtaskID}delete", handels.DeleteSubtask).Methods("DELETE")
	//s.HandleFunc("/tasks/{taskID}/submissions").Methods("GET")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "token"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler := c.Handler(r)
	fmt.Println("Started serving")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill, syscall.SIGKILL)
	go func() {
		err = http.ListenAndServeTLS(":8080", config.CertificateDir, config.PrivateKeyDir, handler)
		if err != nil {
			log.Panic(err)
		}
	}()

	fmt.Println("Run after listen and serve")
	//End

	<-sc
	fmt.Println("Run after SIGKILL or SIGTERM")
}
