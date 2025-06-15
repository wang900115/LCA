package entities

type Channel struct {
	Name  string `json:"name"`
	Users []User `json:"users"`
}
