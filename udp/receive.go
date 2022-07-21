package udp

import (
	"github.com/chuccp/utils/udp/wire"
	"time"
)

type receivePacket struct {
	receiveTime *time.Time
	header *wire.Header
	oob []byte
}

func NewReceivePacket(receiveTime *time.Time,oob []byte,header *wire.Header) *receivePacket {
	return &receivePacket{receiveTime:receiveTime,oob:oob,header:header}
}
