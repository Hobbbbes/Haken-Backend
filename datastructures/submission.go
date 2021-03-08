package datastructures

type Language struct {
	Name          string `json:"name" yaml:"name"`
	Abbreviation  string `json:"abbreviation" yaml:"abbreviation"`
	PreLaunchTask string `json:"-" yaml:"preLaunchTask"`
	LaunchTask    string `json:"-" yaml:"launchTask"`
}

type LanguageName string

//Submission of a User for a given Task
type Submission struct {
	ID     uint64       `json:"id"`
	Author *User        `json:"author"`
	T      *Task        `json:"task"`
	Lang   LanguageName `json:"language"`
	//Source Code path is calculated by ID
}
