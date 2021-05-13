package database

import (
	"log"
	"strings"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//GetTask returns task by id and error
func GetTask(taskID int) (datastructures.Task, error) {
	var task datastructures.Task
	err := db.QueryRow("SELECT id,Name,Author,Description,Group_id FROM Tasks WHERE id = ?", taskID).Scan(
		&task.ID, &task.Name, &task.Author, &task.Description, &task.GroupID)
	return task, err
}

//IsUserAllowedToAccessTask checks if a user is allowed to access a given task
func IsUserAllowedToAccessTask(token string, taskID int) (bool, error) {
	task, err := GetTask(taskID)
	if err != nil {
		log.Println("IsUserAllowedToAccessTask: " + err.Error())
		return false, err
	}
	i, err := isUserInGroup(token, task.GroupID)
	if err != nil {
		log.Println("IsUserAllowedToAccessTask: " + err.Error())
		return false, err
	}
	return i, nil
}

//GetSubtasksForTask returns all subtasks for a task
func GetSubtasksForTask(taskID int, token string) ([]datastructures.Subtask, error) {

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

//AddTask adds a task to a group and returns the fully populated task and an error
func AddTask(task datastructures.Task) (datastructures.Task, error) {
	res, err := db.Exec("INSERT INTO Tasks(Name,Author,Description,Group_id) VALUES (?,?,?,?)",
		task.Name, task.Author, task.Description, task.GroupID)
	if err != nil {
		log.Printf("AddTask:" + err.Error())
		return task, err
	}
	id, err := res.LastInsertId()
	task.ID = int(id)
	if err != nil {
		log.Printf("AddTask:" + err.Error())
		return task, err
	}
	return task, nil

}

func GetTasksForGroups(groupIDs []interface{}) ([]datastructures.Task, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	taskRows, err := db.Query("SELECT id,Name,Author,Description,Group_id FROM `Tasks` WHERE Group_id in (?"+strings.Repeat(",?", len(groupIDs)-1)+")", groupIDs...)
	if err != nil {
		log.Println("GetTasksForGroups: ", err.Error())
		return nil, err
	}
	tasks := make([]datastructures.Task, 0, 20)
	for taskRows.Next() {
		var t datastructures.Task
		err := taskRows.Scan(&t.ID, &t.Name, &t.Author, &t.Description, &t.GroupID)
		if err != nil {
			log.Println("GetTasksForGroups: ", err.Error())
			return nil, err
		}
		t.Author, err = GetUserNameFromToken(t.Author)
		if err != nil {
			log.Println("GetUserNameFromToken: ", err.Error())
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

//AddSubtask adds subtask to database and returns subtask with id
func AddSubtask(t datastructures.Subtask) (datastructures.Subtask, error) {
	res, err := db.Exec("INSERT INTO Subtasks(Points,Name,Tasks_id) VALUES (?,?,?)",
		t.Points, t.Name, t.TaskID)
	if err != nil {
		log.Printf("AddSubtask:" + err.Error())
		return t, err
	}
	id, err := res.LastInsertId()
	t.ID = int(id)
	if err != nil {
		log.Printf("AddSubtask:" + err.Error())
		return t, err
	}
	return t, nil
}

func IsUserAuthorOfTask(token string, taskID int) (bool, error) {
	t, err := GetTask(taskID)
	if err != nil {
		return false, err
	}
	return t.Author == token, nil
}
