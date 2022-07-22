package wire

import (
	"github.com/chuccp/utils/udp/util"
)

type PacketType uint8

const (
	PacketTypeInitial PacketType = iota
	PacketTypeZeroRTT
	PacketTypeHandshake
	PacketTypeRetry
	PacketTypeVersionNegotiation
	PacketTypeOneRTT
)

type Header struct {
	IsLongHeader            bool
	PacketType              PacketType
	Version                 util.VersionNumber
	DestinationConnectionId []byte
	SourceConnectionId      []byte
	FixedBit                bool
	ReservedBits            uint8
	PacketNumberLength      uint8
	Token                   []byte
	RetryToken              []byte
	RetryIntegrityTag       []byte
	ParsedLen uint32
	Length uint32
	PacketNumber            util.PacketNumber
	PacketPayload           []byte
}

func (header *Header) GetPacketNumberLength() uint8 {
	return header.PacketNumberLength
}
func (header *Header) GetFirstByte() byte {
	var b byte = 0
	if header.IsLongHeader {
		b = b | 0x80
	}
	if header.FixedBit {
		b = b | 0x40
	}
	b = b | (uint8(header.PacketType) << 4)
	b = b | (header.ReservedBits << 2)
	b = b | (header.PacketNumberLength - 1)
	return b
}
func (header *Header) SetFirstByte(b byte) {
	header.IsLongHeader = b&0x80 > 0
	header.FixedBit = b&0x40 > 0
	header.PacketType = PacketType(b & 0x30 >> 4)
	header.ReservedBits = 0x0c >> 2
	header.PacketNumberLength = b&0x03 + 1
}

func (header *Header) Write(packetBuffer *util.WriteBuffer) {
	b := header.GetFirstByte()
	packetBuffer.WriteByte(b)
	packetBuffer.WriteBytes(header.Version.ToBytes())
	packetBuffer.WriteUint8LengthBytes(header.DestinationConnectionId)
	packetBuffer.WriteUint8LengthBytes(header.SourceConnectionId)
	if header.PacketType == PacketTypeInitial {
		packetBuffer.WriteVariableLengthBytes(header.Token)
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(header.GetPacketNumberLength(), uint64(header.PacketNumber))
			wr.WriteBytes(header.PacketPayload)
		})
	}
	if header.PacketType == PacketTypeHandshake || header.PacketType == PacketTypeZeroRTT {
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(header.GetPacketNumberLength(), uint64(header.PacketNumber))
			wr.WriteBytes(header.PacketPayload)
		})
	}
	if header.PacketType == PacketTypeRetry {
		packetBuffer.WriteBytes(header.RetryToken)
		packetBuffer.WriteBytes(header.RetryIntegrityTag)
	}
}
func (header *Header) ParseLongHeader(oob []byte) error {
	packetBuffer := util.NewReadBuffer(oob)
	readByte, err := packetBuffer.ReadByte()
	if err != nil {
		return err
	}
	header.SetFirstByte(readByte)
	u32, err := packetBuffer.Read4U32()
	if err != nil {
		return err
	}
	header.Version = util.VersionNumber(u32)
	_, data, err := packetBuffer.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	header.DestinationConnectionId = data
	_, data, err = packetBuffer.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	header.SourceConnectionId = data
	_, data, err = packetBuffer.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	header.Token = data
	header.Length, data, err = packetBuffer.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	header.ParsedLen = uint32(packetBuffer.Offset())
	return nil
}
func (header *Header) ParseExLongHeader(oob []byte)  error {
	packetBuffer := util.NewReadBuffer(oob[header.ParsedLen:header.ParsedLen+header.Length])
	u32, err := packetBuffer.ReadU8LengthU32(header.PacketNumberLength)
	if err!=nil{
		return err
	}
	header.PacketNumber = util.PacketNumber(u32)
	header.PacketPayload,err = packetBuffer.ReadU32Bytes(header.Length-uint32(header.PacketNumberLength))
	return err
}