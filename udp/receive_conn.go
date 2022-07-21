package udp

import (
	"github.com/chuccp/utils/udp/wire"
)


type ReceiveConn struct {
	handshakeComplete     bool
	serverHandshake *wire.ServerHandshake

}

func (rc *ReceiveConn) Read(data []byte)(n int, err error)  {
	return 0,nil
}
func (rc *ReceiveConn) Write(p []byte) (n int, err error) {
	return 0,nil
}
func (rc *ReceiveConn) push(packet *wire.ReceivePacket)  {
	if !rc.handshakeComplete{
		rc.serverHandshake.Handle(packet)
		if rc.serverHandshake.HandshakeStatus==wire.FinishHandshake{
			rc.handshakeComplete = true
		}
	}
}
func (rc *ReceiveConn)handleSinglePacket(receivePacket *wire.ReceivePacket,header *wire.Header)error{
	ParseHeader(receivePacket.Data,header)

	return nil
}

func newReceiveConn() *ReceiveConn {
	return &ReceiveConn{handshakeComplete:false,serverHandshake: wire.NewServerHandshake()}
}