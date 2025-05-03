package event

type EventType string

const (
	EventJoin    EventType = "join"
	EventMessage EventType = "message"
	EventLeave   EventType = "leave"
)
