package entities

type Message struct {
	ID      uint   `json:"id"`
	Sender  uint   `json:"sender"`
	Content string `json:"content"`
}
