package p2p

type RPCtype [1]byte

var (
	INCOMMINGMESSAGE = RPCtype{0x01}
	INCOMMINGSTREAM  = RPCtype{0x02}
	OUTGOINGMESSAGE  = RPCtype{0x03}
	OUTGOINGSTREAM   = RPCtype{0x04}
)

type RPCContext struct {
	From string
	Msg  interface{} // can be either *Message or *Stream
	Sig  []byte
}
