package entities

type Message struct {
	Channel string `json:"channel"`
	User    string `json:"user"`
	Content string `json:"content"`
}
