package database

import (
	"errors"
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

//GetTasksForGroup returns all Tasks a Group contains and checks if requesting user has access
func GetTasksForGroup(userToken string, groupID int) ([]datastructures.Task, error) {
	var necessaryUserToken string
	rows, err := db.Query("SELECT User_Token FROM Group_has_Users WHERE Group_id = ?", groupID)
	if err != nil {
		log.Println("GetTasksForUser: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	var userFound bool = false

	for rows.Next() {
		err := rows.Scan(&necessaryUserToken)
		if err != nil {
			log.Println("GetTasksForUser: " + err.Error())
			return nil, err
		}
		if necessaryUserToken == userToken {
			userFound = true
			break
		}
	}

	if !userFound {
		err = errors.New("User not allowed to view Group details")
		log.Println("GetTasksForUser: " + err.Error())
		return nil, err
	}

	tasks := make([]datastructures.Task, 0, 20)
	taskRows, err := db.Query("SELECT id,Name,Author,Description FROM Tasks WHERE Group_id = ?", groupID)
	if err != nil {
		log.Println("GetTasksForUser: " + err.Error())
		return nil, err
	}
	defer taskRows.Close()
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

func GetGroupsForUser(token string) ([]datastructures.Group, error) {

	groupIDs := make([]interface{}, 0, 20)

	groupIDsRows, err := db.Query("SELECT Group_id FROM Group_has_Users WHERE User_Token = ?", token)

	if err != nil {
		log.Println("GetGroupsForUser: " + err.Error())
		return nil, err
	}
	defer groupIDsRows.Close()

	for groupIDsRows.Next() {
		var groupID int
		err := groupIDsRows.Scan(&groupID)
		if err != nil {
			log.Println("GetGroupsForUser: " + err.Error())
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
	}
	groups := make([]datastructures.Group, 0, 20)
	if len(groupIDs) == 0 {
		return make([]datastructures.Group, 0), nil
	}
	stmt := "SELECT id,Name,Description FROM `Group` WHERE id in (?" + strings.Repeat(",?", len(groupIDs)-1) + ")"
	groupRows, err := db.Query(stmt, groupIDs...)
	if err != nil {
		log.Println("GetGroupsForUser: " + err.Error())
		return nil, err
	}
	for groupRows.Next() {
		var group datastructures.Group
		err := groupRows.Scan(&group.ID, &group.Name, &group.Description)
		if err != nil {
			log.Println("GetGroupsForUser: " + err.Error())
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
