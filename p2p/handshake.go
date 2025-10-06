package p2p

type HandShakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
