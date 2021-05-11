package datastructures

//User stores only access token for now
type User struct {
	Token    string `json:"-"`
	UserName string `json:"name"`
	PwdHash  []byte `json:"-"`
}

type UserLogin struct {
	UserName string `json:"name"`
	Pwd      string `json:"pwd"`
}
