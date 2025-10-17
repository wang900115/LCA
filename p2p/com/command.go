package common

type Command byte

const (
	HEARTBEAT      Command = 0x00
	PEERINFO       Command = 0x01
	PEERACK        Command = 0x02
	MESSAGESEND    Command = 0x04
	MESSAGESENDACK Command = 0x05
)

type Message byte

const (
	PUBLIC  Message = 0x00
	PRIVATE Message = 0x01
)
