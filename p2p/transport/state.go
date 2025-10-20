package transport

import (
	"sync"
	"sync/atomic"

	"github.com/wang900115/LCA/p2p"
)

type stateError struct{ msg string }

func (e *stateError) Error() string { return e.msg }

var (
	errExceedOutBoundLimit      = &stateError{"exceed outbound peer limit"}
	errExceedInBoundLimit       = &stateError{"exceed inbound peer limit"}
	errPeerNotFoundInboundKeys  = &stateError{"peer not found in inbound keys"}
	errPeerNotFoundOutboundKeys = &stateError{"peer not found in outbound keys"}
)

type state struct {
	mu sync.RWMutex

	outBoundCnt  atomic.Int32
	inBoundCnt   atomic.Int32
	handshakeCnt atomic.Int32

	outBoundLi    int
	inBoundLi     int
	outBoundPeers map[string]p2p.Peer
	inBoundPeers  map[string]p2p.Peer
	outKeys       map[string][]byte
	inKeys        map[string][]byte
}

func NewState(outBoundLimit, inBoundLimit int) *state {
	return &state{
		outBoundLi:    outBoundLimit,
		inBoundLi:     inBoundLimit,
		outBoundPeers: make(map[string]p2p.Peer),
		inBoundPeers:  make(map[string]p2p.Peer),
	}
}

// Count returns the current counts of outBound, inBound, and handshake peers.
func (s *state) Count() (outBound, inBound, handshake int) {
	return int(s.outBoundCnt.Load()), int(s.inBoundCnt.Load()), int(s.handshakeCnt.Load())
}

// Limit returns the limits for outBound and inBound peers.
func (s *state) Limit() (outBoundLi, inBoundLi int) {
	return s.outBoundLi, s.inBoundLi
}

// OutPeers returns the map of outbound peers.
func (s *state) OutPeers() map[string]p2p.Peer {
	return s.outBoundPeers
}

// InPeers returns the map of inbound peers.
func (s *state) InPeers() map[string]p2p.Peer {
	return s.inBoundPeers
}

// Increment / Decrement counters using atomic operations
func (s *state) IncOutBound()  { s.outBoundCnt.Add(1) }
func (s *state) DecOutBound()  { s.outBoundCnt.Add(-1) }
func (s *state) IncInBound()   { s.inBoundCnt.Add(1) }
func (s *state) DecInBound()   { s.inBoundCnt.Add(-1) }
func (s *state) IncHandShake() { s.handshakeCnt.Add(1) }
func (s *state) DecHandShake() { s.handshakeCnt.Add(-1) }

// AddOutPeer adds a peer to the outbound peer map.
func (s *state) AddOutPeer(peer p2p.Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if int(s.outBoundCnt.Load()) > s.outBoundLi {
		return errExceedOutBoundLimit
	}
	s.outBoundPeers[peer.Addr()] = peer
	return nil
}

// RemoveOutPeer removes a peer from the outbound peer map.
func (s *state) RemoveOutPeer(peer p2p.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.outBoundPeers, peer.Addr())
}

// AddInPeer adds a peer to the inbound peer map.
func (s *state) AddInPeer(peer p2p.Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if int(s.inBoundCnt.Load()) > s.inBoundLi {
		return errExceedInBoundLimit
	}
	s.inBoundPeers[peer.Addr()] = peer
	return nil
}

// RemoveInPeer removes a peer from the inbound peer map.
func (s *state) RemoveInPeer(peer p2p.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.inBoundPeers, peer.Addr())
}

// SetOutKey sets the outbound encryption key for a given peer address.
func (s *state) SetOutKey(addr string, key []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.outKeys[addr] = key
}

// SetInKey sets the inbound encryption key for a given peer address.
func (s *state) SetInKey(addr string, key []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inKeys[addr] = key
}

// GetOutKey retrieves the outbound encryption key for a given peer address.
func (s *state) GetOutKey(addr string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, exist := s.outKeys[addr]
	if !exist {
		return nil, errPeerNotFoundOutboundKeys
	}
	return key, nil
}

// GetInKey retrieves the inbound encryption key for a given peer address.
func (s *state) GetInKey(addr string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, exist := s.inKeys[addr]
	if !exist {
		return nil, errPeerNotFoundInboundKeys
	}
	return key, nil
}
