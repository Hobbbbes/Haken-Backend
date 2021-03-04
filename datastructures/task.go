package datastructures

//Task describes all important information
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	GroupID     int    `json:"groupID"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsAdmin     bool   `json:"isadmin"`
}

type GroupWithTasks struct {
	Group
	Tasks []Task `json:"tasks"`
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
