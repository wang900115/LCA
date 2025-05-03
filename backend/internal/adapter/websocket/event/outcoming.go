package event

type OutcomingEvent struct {
	Type     EventType `json:"type"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
}
