package event

const (
	EventJoin    = "join"
	EventMessage = "message"
	EventLeave   = "leave"
)

type MessagePayload struct {
	UUID      string `json:"uuid"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
