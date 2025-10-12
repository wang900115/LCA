package packet

// 1 byte for type
// 1 bytes for command
// 4 bytes for length
// n bytes for payload
// 1 + 1 +4 = 6 bytes header
// 4 bytes footer (checksum)

const (
	INCOMMINGMESSAGE = 0x1
	INCOMMINGSTREAM  = 0x2
	OUTGOINGMESSAGE  = 0x3
	OUTGOINGSTREAM   = 0x4
)

const (
	HEARTBEAT      = 0x00
	PEERINFO       = 0x01
	PEERACK        = 0x02
	PEERERROR      = 0x03
	BETCREATE      = 0x03
	BETACK         = 0x04
	BETERROR       = 0x05
	RESETTLECREATE = 0x06
	RESETTLEACK    = 0x07
	RESETTLEERROR  = 0x08
	ROUNDSTART     = 0x09
	ROUNDEND       = 0x0A
	ROUNDWAIT      = 0x0B
)

type RPC struct {
	From string
	Msg  interface{} // can be either *Message or *Stream
	Sig  []byte
}
