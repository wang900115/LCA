package p2p

type RPCContext struct {
	From string
	Msg  []byte // is will coming to message struct
	Sig  []byte
}

// GetFrom returns the sender of the RPC
func (r *RPCContext) GetFrom() string {
	return r.From
}

// getMsg returns the message of the RPC
func (r *RPCContext) getMsg() interface{} {
	return r.Msg
}

// getSig returns the signature of the RPC
func (r *RPCContext) getSig() []byte {
	return r.Sig
}

func (r *RPCContext) Encode(msg interface{}) ([]byte, error) {

	return nil, nil
}

// DecodeRPC2MSG decodes a byte slice into a message
func DecodeRPC2MSG(data []byte) (Message, error) {
	var msg Message
	return msg, nil
}

func ValidateRPC(rpc interface{}) error {

	return nil
}

// Call performs an RPC call to the given method with args and fills the reply
func (r *RPCContext) Call(method string, args interface{}, reply interface{}) error {
	return nil
}
