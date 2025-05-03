package entities

type User struct {
	UUID        string `json:"uuid"`
	Username    string `json:"username"`
	ChannelUUID string `json:"channel_uuid"`
	Status      string `json:"status"`
}
