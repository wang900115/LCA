package network

type RPCContent struct {
	From string
	Msg  []byte // is will coming to message struct
	Sig  []byte
}
