package database

import (
	"log"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

//AddSubmission adds submission to the database and returns submission with id
func AddSubmission(sub datastructures.Submission) (datastructures.Submission, error) {
	res, err := db.Exec("INSERT INTO Submission(User_Token,Tasks_id,Tasks_Group_id) VALUES (?,?,?)",
		sub.Author, sub.TaskID, sub.GroupID)
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
		res.Sub.ID, res.Sub.Author, res.Sub.TaskID, res.Sub.GroupID, res.Subt.ID, res.Success)
	if err != nil {
		log.Printf("AddResult:" + err.Error())

	}
	return err
}
