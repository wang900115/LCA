package p2p

type RPC struct {
	From string
	Msg  interface{} // can be either *Message or *Stream
	Sig  []byte
}
