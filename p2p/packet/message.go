package packet

const (
	INCOMMINGMESSAGE = 0x1
	INCOMMINGSTREAM  = 0x2
)

type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
