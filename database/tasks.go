package database

import (
	"log"
	"strings"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

func getTask(taskID int) (datastructures.Task, error) {
	var task datastructures.Task
	err := db.QueryRow("SELECT id,Name,Author,Description FROM Tasks WHERE id = ?", taskID).Scan(
		&task.ID, &task.Name, &task.Author, &task.Description)
	return task, err
}

//GetTasksForUser returns all Tasks a User has access to
func GetTasksForUser(userToken string) ([]datastructures.Task, error) {
	var user datastructures.User
	err := db.QueryRow("SELECT id,Token FROM Users WHERE Token = ?", userToken).Scan(
		&user.ID, &user.Token)
	if err != nil {
		log.Println("GetTasksForUser: " + err.Error())
		return nil, err
	}
	taskIDs := make([]interface{}, 0, 20)

	taskIDsRows, err := db.Query("SELECT Tasks_id FROM Tasks_has_Users WHERE Users_id = ?", user.ID)

	if err != nil {
		log.Println("GetTasksForUser: " + err.Error())
		return nil, err
	}
	for taskIDsRows.Next() {
		var taskID int
		err := taskIDsRows.Scan(&taskID)
		if err != nil {
			log.Println("GetTasksForUser: " + err.Error())
			return nil, err
		}
		taskIDs = append(taskIDs, taskID)
	}
	tasks := make([]datastructures.Task, 0, 20)

	stmt := "SELECT id,Name,Author,Description FROM Tasks WHERE id in (?" + strings.Repeat(",?", len(taskIDs)-1) + ")"
	taskRows, err := db.Query(stmt, taskIDs...)
	for taskRows.Next() {
		var task datastructures.Task
		err := taskRows.Scan(&task.ID, &task.Name, &task.Author, &task.Description)
		if err != nil {
			log.Println("GetTasksForUser: " + err.Error())
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
