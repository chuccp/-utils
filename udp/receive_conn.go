package udp

import (
	"github.com/chuccp/utils/udp/wire"
)


type ReceiveConn struct {
	handshakeComplete     bool
}

func (rc *ReceiveConn) Read(data []byte)(n int, err error)  {
	return 0,nil
}
func (rc *ReceiveConn) Write(p []byte) (n int, err error) {
	return 0,nil
}
func (rc *ReceiveConn) push(packet *receivePacket)  {






}
func (rc *ReceiveConn)handleSinglePacket(receivePacket *receivePacket,header *wire.Header)error{
	ParseHeader(receivePacket.oob,header)

	return nil
}

func newReceiveConn() *ReceiveConn {
	return &ReceiveConn{handshakeComplete:false}
}