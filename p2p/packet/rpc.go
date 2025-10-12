package packet

const (
	INCOMMINGMESSAGE = 0x1
	INCOMMINGSTREAM  = 0x2
)

type RPC struct {
	From    string
	Content Message
	Stream  bool
}

type Message struct {
	To        string
	Payload   []byte
	Timestamp int64
}
