package udp

import "time"

type receivePacket struct {
	receiveTime *time.Time
	oob []byte
}

func NewReceivePacket(receiveTime *time.Time,oob []byte) *receivePacket {
	return &receivePacket{receiveTime:receiveTime,oob:oob}
}
