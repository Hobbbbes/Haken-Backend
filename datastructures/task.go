package datastructures

//Task describes all important information
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	GroupID     int    `json:"-"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

//Subtask holds Information about one run of the programm
type Subtask struct {
	ID     int    `json:"id"`
	Points int    `json:"points"`
	Name   string `json:"name"`
	TaskID int    `json:"-"`
	//In and output file paths are calculated on the fly by ID
}

//Result connects a submission with a subtask
type Result struct {
	Submiss    *Submission
	Subt       *Subtask
	Points     int
	ResultCode string
}

//Submission of a User for a given Task
type Submission struct {
	ID     uint64
	Author *User
	T      *Task
	//Source Code path is calculated by ID
}
