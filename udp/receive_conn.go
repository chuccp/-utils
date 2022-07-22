package udp

import (
	"github.com/chuccp/utils/udp/wire"
)

type ReceiveConn struct {
	handshakeComplete bool
	serverHandshake   *wire.ServerHandshake
}

func (rc *ReceiveConn) Read(data []byte) (n int, err error) {
	return 0, nil
}
func (rc *ReceiveConn) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (rc *ReceiveConn) push(packet *wire.ReceivePacket) {
	rc.handleSinglePacket(packet)
}
func (rc *ReceiveConn) handleSinglePacket(receivePacket *wire.ReceivePacket) error {
	if receivePacket.Header.IsLongHeader {
		rc.serverHandshake.Handle(receivePacket)
	}
	return nil
}

func newReceiveConn() *ReceiveConn {
	return &ReceiveConn{handshakeComplete: false, serverHandshake: wire.NewServerHandshake()}
}
