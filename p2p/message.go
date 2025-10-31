package p2p

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"sync/atomic"
// 	"time"

// 	"github.com/ethereum/go-ethereum/event"
// 	"github.com/ethereum/go-ethereum/p2p/enode"
// )

// type Msg struct {
// 	Code       uint64
// 	Size       uint32
// 	Payload    io.Reader
// 	ReceivedAt time.Time

// 	meterCap  Cap
// 	meterCode uint64
// 	meterSize uint32
// }

// func (msg Msg) Decode(val interface{}) error {
// 	data, err := io.ReadAll(msg.Payload)
// 	if err != nil {
// 		return err
// 	}
// 	buf := bytes.NewBuffer(data)
// 	return binary.Read(buf, binary.BigEndian, val)
// }

// func (msg Msg) String() string {
// 	return fmt.Sprintf("msg #%v (%v bytes)", msg.Code, msg.Size)
// }

// func (msg Msg) Discard() error {
// 	_, err := io.Copy(io.Discard, msg.Payload)
// 	return err
// }

// func (msg Msg) Time() time.Time {
// 	return msg.ReceivedAt
// }

// type MsgReader interface {
// 	ReadMsg() (Msg, error)
// }

// type MsgWriter interface {
// 	WriteMsg(Msg) error
// }

// type MsgReadWriter interface {
// 	MsgReader
// 	MsgWriter
// }

// func Send(w MsgWriter, msgCode uint64, data interface{}) error {
// 	var buf bytes.Buffer
// 	if err := binary.Write(&buf, binary.BigEndian, data); err != nil {
// 		return err
// 	}
// 	return w.WriteMsg(Msg{
// 		Code:    msgCode,
// 		Size:    uint32(buf.Len()),
// 		Payload: &buf,
// 	})
// }

// type eofSignal struct {
// 	wrapped io.Reader
// 	count   uint32
// 	eof     chan<- struct{}
// }

// func (r *eofSignal) Read(buf []byte) (n int, err error) {
// 	if r.count == 0 {
// 		if r.eof != nil {
// 			r.eof <- struct{}{}
// 			r.eof = nil
// 		}
// 		return 0, io.EOF
// 	}

// 	max := len(buf)
// 	if int(r.count) < len(buf) {
// 		max = int(r.count)
// 	}
// 	n, err = r.wrapped.Read(buf[:max])
// 	r.count -= uint32(n)
// 	if (err != nil || r.count == 0) && r.eof != nil {
// 		r.eof <- struct{}{}
// 		r.eof = nil
// 	}
// 	return n, err
// }

// func MsgPipe() (*MsgPipeRW, *MsgPipeRW) {
// 	var (
// 		c1, c2  = make(chan Msg), make(chan Msg)
// 		closing = make(chan struct{})
// 		closed  = new(atomic.Bool)
// 		rw1     = &MsgPipeRW{c1, c2, closing, closed}
// 		rw2     = &MsgPipeRW{c2, c1, closing, closed}
// 	)
// 	return rw1, rw2
// }

// var ErrPipeClosed = errors.New("p2p: read or write on closed message pipe")

// type MsgPipeRW struct {
// 	w       chan<- Msg
// 	r       <-chan Msg
// 	closing chan struct{}
// 	closed  *atomic.Bool
// }

// func (p *MsgPipeRW) Close() error {
// 	if p.closed.Swap(true) {
// 		return nil
// 	}
// 	close(p.closing)
// 	return nil
// }

// func (p *MsgPipeRW) WriteMsg(msg Msg) error {
// 	if !p.closed.Load() {
// 		consumed := make(chan struct{}, 1)
// 		msg.Payload = &eofSignal{msg.Payload, msg.Size, consumed}
// 		select {
// 		case p.w <- msg:
// 			if msg.Size > 0 {
// 				select {
// 				case <-consumed:
// 				case <-p.closing:
// 				}
// 			}
// 			return nil
// 		case <-p.closing:
// 		}
// 	}
// 	return ErrPipeClosed
// }

// func (p *MsgPipeRW) ReadMsg() (Msg, error) {
// 	if !p.closed.Load() {
// 		select {
// 		case msg := <-p.r:
// 			return msg, nil
// 		case <-p.closing:
// 		}
// 	}
// 	return Msg{}, ErrPipeClosed
// }

// func ExpectMsg(r MsgReader, code uint64, content interface{}) error {
// 	msg, err := r.ReadMsg()
// 	if err != nil {
// 		return err
// 	}
// 	if msg.Code != code {
// 		return fmt.Errorf("expected msg code %d, got %d", code, msg.Code)
// 	}
// 	data, err := io.ReadAll(msg.Payload)
// 	if err != nil {
// 		return err
// 	}
// 	expected := new(bytes.Buffer)
// 	if err := binary.Write(expected, binary.BigEndian, content); err != nil {
// 		return err
// 	}
// 	if !bytes.Equal(data, expected.Bytes()) {
// 		return fmt.Errorf("message content does not match expected content")
// 	}
// 	return nil
// }

// type msgEventer struct {
// 	MsgReadWriter

// 	feed          *event.Feed
// 	peerID        enode.ID
// 	Protocol      string
// 	localAddress  string
// 	remoteAddress string
// }

// func newMsgEventer(rw MsgReadWriter, feed *event.Feed, peerID enode.ID, proto, remote, local string) *msgEventer {
// 	return &msgEventer{
// 		MsgReadWriter: rw,
// 		feed:          feed,
// 		peerID:        peerID,
// 		Protocol:      proto,
// 		remoteAddress: remote,
// 		localAddress:  local,
// 	}
// }

// func (ev *msgEventer) ReadMsg() (Msg, error) {
// 	msg, err := ev.MsgReadWriter.ReadMsg()
// 	if err != nil {
// 		return msg, err
// 	}
// 	ev.feed.Send(&PeerEvent{
// 		Type:          PeerEventTypeMsgRecv,
// 		Peer:          ev.peerID,
// 		Protocol:      ev.Protocol,
// 		MsgCode:       &msg.Code,
// 		MsgSize:       &msg.Size,
// 		LocalAddress:  ev.localAddress,
// 		RemoteAddress: ev.remoteAddress,
// 	})
// 	return msg, nil
// }

// func (ev *msgEventer) WriteMsg(msg Msg) error {
// 	err := ev.MsgReadWriter.WriteMsg(msg)
// 	if err != nil {
// 		return err
// 	}
// 	ev.feed.Send(&PeerEvent{
// 		Type:          PeerEventTypeMsgSend,
// 		Peer:          ev.peerID,
// 		Protocol:      ev.Protocol,
// 		MsgCode:       &msg.Code,
// 		MsgSize:       &msg.Size,
// 		LocalAddress:  ev.localAddress,
// 		RemoteAddress: ev.remoteAddress,
// 	})
// 	return nil
// }

// func (ev *msgEventer) Close() error {
// 	if closer, ok := ev.MsgReadWriter.(io.Closer); ok {
// 		return closer.Close()
// 	}
// 	return nil
// }
