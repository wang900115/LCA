package entities

type Channel struct {
	UUID  string `json:"uuid"`
	Users []User `json:"users"`
}
