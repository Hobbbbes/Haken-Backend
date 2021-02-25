package main

import (
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
	database.InitDB(config.DBName, config.DBUser, config.DBPassword)
	r := mux.NewRouter().StrictSlash(true)
	//Use for unautherized route

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(handels.AuthToken)

	s.HandleFunc("/tasks", handels.GetTasks).Methods("GET")
	s.HandleFunc("/task/{id}", handels.GetTask).Methods("GET")

	http.ListenAndServe(":8080", s)
	//End
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
