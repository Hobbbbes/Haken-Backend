package datastructures

type Language struct {
	Name          string `json:"name" yaml:"name"`
	Abbreviation  string `json:"abbreviation" yaml:"abbreviation"`
	PreLaunchTask string `json:"-" yaml:"preLaunchTask"`
	LaunchTask    string `json:"-" yaml:"launchTask"`
}

//Submission of a User for a given Task
type Submission struct {
	ID     uint64
	Author *User
	T      *Task
	Lang   Language
	//Source Code path is calculated by ID
}
