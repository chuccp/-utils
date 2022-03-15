package udp

import (
	"io"
)

type PacketType uint8

const (
	PacketTypeInitial PacketType = iota
	PacketTypeRetry
	PacketTypeHandshake
	PacketType0RTT
)

func ParseConnectionID(data []byte) ([]byte, error) {
	isLongHeader := data[0]&0x80 > 0
	if isLongHeader {
		destConnIDLen := int(data[5])
		connectionId := data[6 : 6+destConnIDLen]
		return connectionId, nil
	}
	return nil, io.EOF
}

type Header struct {
	typeByte        byte
	packetType      PacketType
	longHeader      bool
	version         []byte
	desConnIdLen    int
	desConnId       []byte
	sourceConnIdLen int
	sourceConnId    []byte
	tokenLen        int
	token           []byte
	length          uint64
	parsedLen       int
	payload         []byte

}

func parseLongHeader(data []byte) *Header {
	header := &Header{}
	header.typeByte = data[0]
	if header.typeByte&128>0{
		header.longHeader = true
	}
	header.packetType = PacketType((header.typeByte & 0x30) >> 4)
	header.version = data[1:5]
	header.desConnIdLen = int(data[5])
	if header.desConnIdLen > 0 {
		header.desConnId = data[6 : 6+header.desConnIdLen]
	}
	header.sourceConnIdLen = int(data[6+header.desConnIdLen])
	if header.sourceConnIdLen > 0 {
		header.sourceConnId = data[7+header.desConnIdLen : 7+header.desConnIdLen+header.sourceConnIdLen]
	}
	header.tokenLen = int(data[7+header.desConnIdLen+header.sourceConnIdLen])
	if header.tokenLen > 0 {
		header.token = data[8+header.desConnIdLen+header.sourceConnIdLen : 8+header.desConnIdLen+header.sourceConnIdLen+header.tokenLen]
	}
	index := 8 + header.desConnIdLen + header.sourceConnIdLen + header.tokenLen
	length, ty := readLength(data[index:])
	header.length = length
	header.parsedLen = index+ty
	header.payload = data[index+ty : index+ty+int(length)]
	return header
}

func readLength(data []byte) (uint64, int) {
	lenType := data[0]
	switch lenType >> 6 {
	case 0:
		return uint64(data[0] & 63), 1
	case 1:
		return uint64(data[0]&63)<<8 | uint64(data[1]), 2
	case 2:
		return uint64(data[0]&63)<<16 | uint64(data[1])<<8 | uint64(data[2]), 3
	default:
		return uint64(data[0]&63)<<24 | uint64(data[1])<<16 | uint64(data[2])<<8 | uint64(data[3]), 4
	}
}
