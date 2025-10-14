package common

type Command byte

const (
	HEARTBEAT         Command = 0x00
	PEERINFO          Command = 0x01
	PEERACK           Command = 0x02
	PEERERROR         Command = 0x03
	TRANSACTIONCREATE Command = 0x04
	TRANSACTIONACK    Command = 0x05
	TRANSACTIONERROR  Command = 0x06
	BLOCKCREATE       Command = 0x07
	BLOCKACK          Command = 0x08
	BLOCKERROR        Command = 0x09
)

type RPCtype byte

const (
	JOINCHANNELMESSAGE  RPCtype = 0x01
	LEAVECHANNELMESSAGE RPCtype = 0x02
)
