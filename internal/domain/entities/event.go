package entities

type Event struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	EventType string  `json:"type"`
	StartTime int64   `json:"startTime"`
	EndTime   int64   `json:"endTime"`
	Users     []*User `json:"users"`
}
