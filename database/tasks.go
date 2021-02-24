package database

import "github.com/poodlenoodle42/Hacken-Backend/datastructures"

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
		return nil, err
	}
	taskIDs := make([]int, 0, 20)

	taskIDsRows, err := db.Query("SELECT Tasks_id FROM Tasks_has_Users WHERE User_id = ?", user.ID)

	if err != nil {
		return nil, err
	}
	for taskIDsRows.Next() {
		var taskID int
		err := taskIDsRows.Scan(&taskID)
		if err != nil {
			return nil, err
		}
		taskIDs = append(taskIDs, taskID)
	}
	tasks := make([]datastructures.Task, 0, 20)
	for id := range taskIDs {
		task, err := getTask(id)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
