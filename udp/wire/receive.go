package wire

import (
	"time"
)

type ReceivePacket struct {
	ReceiveTime *time.Time
	Header      *Header
	Data        []byte
}

func NewReceivePacket(receiveTime *time.Time,data []byte,header *Header) *ReceivePacket {
	return &ReceivePacket{ReceiveTime: receiveTime, Data:data,Header:header}
}
