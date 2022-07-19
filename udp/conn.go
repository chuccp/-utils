package udp

import (
	"github.com/chuccp/utils/udp/wire"
	"time"
)

type RawConn struct {

}

func (rc *RawConn) Read(data []byte)(n int, err error)  {
	return 0,nil
}
func (rc *RawConn) Write(p []byte) (n int, err error) {
	return 0,nil
}
func (rc *RawConn) push(p []byte) (n int, err error) {
	nowDate:=time.Now()
	receivePacket:=NewReceivePacket(&nowDate,p)
	var header wire.Header
	err = rc.handleSinglePacket(receivePacket, &header)
	if err != nil {
		return 0, err
	}
	return 0,nil
}
func (rc *RawConn)handleSinglePacket(receivePacket *receivePacket,header *wire.Header)error{
	ParseHeader(receivePacket.oob,header)

	return nil
}

func newRawConn() *RawConn {
	return &RawConn{}
}