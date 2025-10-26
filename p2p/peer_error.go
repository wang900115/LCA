package p2p

import (
	"errors"
	"fmt"
)

const (
	errInvalidMsgCode = iota
	errInvalidMsg
)

var error2String = map[int]string{
	errInvalidMsgCode: "invalid message code",
	errInvalidMsg:     "invalid message",
}

type peerError struct {
	code    int
	message string
}

func newPeerError(code int, format string, v ...interface{}) *peerError {
	desc, ok := error2String[code]
	if !ok {
		panic("invalid error code ")
	}
	err := &peerError{
		code:    code,
		message: desc,
	}
	if format != "" {
		err.message += ": " + fmt.Sprintf(format, v...)
	}
	return err
}

func (e *peerError) Error() string {
	return e.message
}

var errProtocolReturned = errors.New("protocol returned")

type DiscReason uint8

const (
	DiscRequested DiscReason = iota
	DiscNetworkError
	DiscProtocolError
	DiscUselessPeer
	DiscTooManyPeers
	DiscAlreadyConnected
	DiscIncompatibleVersion
	DiscInvalidIdentity
	DiscQuitting
	DiscUnexpectedIdentity
	DiscSelf
	DiscReadTimeout
	DiscSubprotocolError = DiscReason(0x10)

	DiscInvalid = 0xff
)

var discReason2String = map[DiscReason]string{
	DiscRequested:           "disconnect request",
	DiscNetworkError:        "network error",
	DiscProtocolError:       "breach of protocol",
	DiscUselessPeer:         "useless peer",
	DiscTooManyPeers:        "too many peers",
	DiscAlreadyConnected:    "already connected",
	DiscIncompatibleVersion: "incompatible p2p protocol version",
	DiscInvalidIdentity:     "invalid node identity",
	DiscQuitting:            "client quitting",
	DiscUnexpectedIdentity:  "unexpected identity",
	DiscSelf:                "connected to self",
	DiscReadTimeout:         "read timeout",
	DiscSubprotocolError:    "subprotocol error",
	DiscInvalid:             "invalid reason",
}

func (d DiscReason) String() string {
	if len(discReason2String) <= int(d) || discReason2String[d] == "" {
		return fmt.Sprintf("unknown reason %d", d)
	}
	return discReason2String[d]
}

func (d DiscReason) Error() string {
	return d.String()
}

func discReansonForError(err error) DiscReason {
	if reason, ok := err.(DiscReason); ok {
		return reason
	}
	if errors.Is(err, errProtocolReturned) {
		return DiscQuitting
	}
	peerError, ok := err.(*peerError)
	if ok {
		switch peerError.code {
		case errInvalidMsgCode, errInvalidMsg:
			return DiscProtocolError
		default:
			return DiscSubprotocolError
		}

	}
	return DiscSubprotocolError
}
