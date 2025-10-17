package p2p

type RPCContent struct {
	From string
	Msg  []byte // is will coming to message struct
	Sig  []byte
}
