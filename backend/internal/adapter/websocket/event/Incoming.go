package event

type IncomingEvent struct {
	Type    EventType `json:"type"`
	Content string    `json:"content"`
}
