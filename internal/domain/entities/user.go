package entities

type User struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Status   string `json:"status"`
}
