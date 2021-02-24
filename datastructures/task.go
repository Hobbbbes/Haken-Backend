package datastructures

//Task describes all important information
type Task struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

//Subtask holds Information about one run of the programm
type Subtask struct {
	ID       uint64
	Points   int
	MainTask *Task
	//In and output file paths are calculated on the fly by ID
}

//TasksPerUser holds all Tasks a given User can make
type TasksPerUser struct {
	Us    *User
	Tasks []Task
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
