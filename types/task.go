package types

//Task describes all important information
type Task struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	DocumentURL string `json:"documentUrl"`
}
