package udp

import (
	"fmt"
	"github.com/chuccp/utils/log"
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
	desConnId       ConnectionID
	sourceConnIdLen int
	sourceConnId    ConnectionID
	tokenLen        int
	token           []byte
	length          uint32
	parsedLen       int
	payload         []byte

}
type ConnectionID []byte
func (c ConnectionID) String() string {
	if c.Len() == 0 {
		return "(empty)"
	}
	return fmt.Sprintf("%x", c.Bytes())
}
func (c ConnectionID) Len() int {
	return len(c)
}
func (c ConnectionID) Bytes() []byte {
	return []byte(c)
}
func parseLongHeader(data []byte) *Header {
	log.Info(data)
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
	length, ty := ReadBytesVariableLength(data[index:])
	log.Info("length:",length)
	header.length = length
	header.parsedLen = index+ty
	header.payload = data[index+ty : index+ty+int(length)]
	log.Info("payload:",header.payload)
	return header
}


