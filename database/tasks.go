package database

import (
	"errors"
	"log"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

func getTask(taskID int) (datastructures.Task, error) {
	var task datastructures.Task
	err := db.QueryRow("SELECT id,Name,Author,Description,Group_id FROM Tasks WHERE id = ?", taskID).Scan(
		&task.ID, &task.Name, &task.Author, &task.Description, &task.GroupID)
	return task, err
}

func GetSubtasksForTasks(taskID int, token string) ([]datastructures.Subtask, error) {
	task, err := getTask(taskID)
	if err != nil {
		log.Println("GetSubtasksForTasks: " + err.Error())
		return nil, err
	}
	group, err := getGroup(task.GroupID)
	if err != nil {
		log.Println("GetSubtasksForTasks: " + err.Error())
		return nil, err
	}
	isUserAllowedToAccess, err := isUserInGroup(token, group.ID)
	if err != nil {
		log.Println("GetSubtasksForTasks: " + err.Error())
		return nil, err
	}
	if !isUserAllowedToAccess {
		err := errors.New("User not allowed to view Group details")
		return nil, err
	}

	subtaskRows, err := db.Query("SELECT id,Points,Tasks_id,Name FROM `Subtasks` WHERE Tasks_id = ?", taskID)
	if err != nil {
		log.Println("GetSubtasksForTasks: " + err.Error())
		return nil, err
	}
	defer subtaskRows.Close()
	subtasks := make([]datastructures.Subtask, 0, 10)
	for subtaskRows.Next() {
		var subtask datastructures.Subtask
		err := subtaskRows.Scan(&subtask.ID, &subtask.Points, &subtask.TaskID, &subtask.Name)
		if err != nil {
			log.Println("GetSubtasksForTasks: " + err.Error())
			return nil, err
		}
		subtasks = append(subtasks, subtask)
	}
	return subtasks, nil
}
