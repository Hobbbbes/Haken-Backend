package database

import (
	"log"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//AddSubmission adds submission to the database and returns submission with id
func AddSubmission(sub datastructures.Submission) (datastructures.Submission, error) {
	res, err := db.Exec("INSERT INTO Submission(LanguageAbbreviation,User_Token,Tasks_id,Tasks_Group_id) VALUES (?,?,?,?)",
		sub.LanguageAbbreviation, sub.Author, sub.TaskID, sub.GroupID)
	if err != nil {
		log.Printf("AddSubmission:" + err.Error())
		return sub, err
	}
	id, err := res.LastInsertId()
	sub.ID = int(id)
	if err != nil {
		log.Printf("AddSubmission:" + err.Error())
		return sub, err
	}
	return sub, nil
}

func AddResult(res datastructures.Result) error {
	_, err := db.Exec("INSERT INTO Result(Submission_id,Submission_User_Token,Submission_Tasks_id,Submission_Tasks_Group_id,Subtasks_id,Success) VALUES (?,?,?,?,?,?)",
		res.Sub.ID, res.Sub.Author, res.Sub.TaskID, res.Sub.GroupID, res.Subt.ID, res.Stat.ExitCode)
	if err != nil {
		log.Printf("AddResult:" + err.Error())

	}
	return err
}

func GetSubmissionsForTask(taskID int) ([]datastructures.Submission, error) {
	submissionRow, err := db.Query("SELECT id,LanguageAbbreviation,User_Token,Tasks_id,Tasks_Group_id FROM `Submission` WHERE Tasks_id = ?", taskID)
	if err != nil {
		log.Println("GetSubmissionsForTask: " + err.Error())
		return nil, err
	}
	defer submissionRow.Close()
	submissions := make([]datastructures.Submission, 0, 10)
	for submissionRow.Next() {
		var submission datastructures.Submission
		err := submissionRow.Scan(&submission.ID, &submission.LanguageAbbreviation, &submission.Author, &submission.TaskID, &submission.GroupID)
		if err != nil {
			log.Println("GetSubmissionsForTask: " + err.Error())
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func GetSubmission(id int) (datastructures.Submission, error) {
	var sub datastructures.Submission
	err := db.QueryRow("SELECT id,User_Token,Tasks_id,Tasks_Group_id FROM Submission WHERE id = ?", id).Scan(
		&sub.ID, &sub.Author, &sub.TaskID, &sub.GroupID)
	return sub, err
}

func GetResultsForSubmission(subID int) ([]datastructures.Status, error) {
	submissionRow, err := db.Query("SELECT id,User_Token,Tasks_id,Tasks_Group_id FROM `Submission` WHERE Tasks_id = ?", subID)
	if err != nil {
		log.Println("GetResultsForSubmission: " + err.Error())
		return nil, err
	}
	defer submissionRow.Close()
	submissions := make([]datastructures.Submission, 0, 10)
	for submissionRow.Next() {
		var submission datastructures.Submission
		err := submissionRow.Scan(&submission.ID, &submission.Author, &submission.TaskID, &submission.GroupID)
		if err != nil {
			log.Println("GetResultsForSubmission: " + err.Error())
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return nil, nil
}
