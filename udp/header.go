package udp

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

type LongHeader struct {
	IsLongHeader bool
	FixedBit     bool

	LongPacketType     packetType
	ReservedBits       uint8
	PacketNumberLength uint8

	Version util.VersionNumber
	//DestinationConnectionIdLength uint8
	DestinationConnectionId []byte

	//SourceConnectionIdLength uint8
	SourceConnectionId []byte

	//TokenVariableLength uint32
	Token []byte
	//LengthVariable      uint32
	PacketNumber  util.PacketNumber
	PacketPayload []byte

	RetryToken        []byte
	RetryIntegrityTag []byte
}

func (longHeader *LongHeader) GetPacketNumberLength() uint8 {
	return longHeader.PacketNumberLength
}
func (longHeader *LongHeader) GetFirstByte() byte {
	var b byte = 0
	if longHeader.IsLongHeader {
		b = b | 0x80
	}
	if longHeader.FixedBit {
		b = b | 0x40
	}
	b = b | (uint8(longHeader.LongPacketType) << 4)
	b = b | (longHeader.ReservedBits << 2)
	b = b | (longHeader.PacketNumberLength - 1)
	return b
}
func (longHeader *LongHeader) SetFirstByte(b byte) {
	longHeader.IsLongHeader = b&0x80 > 0
	longHeader.FixedBit = b&0x40 > 0
	longHeader.LongPacketType = packetType(b & 0x30 >> 4)
	longHeader.ReservedBits = 0x0c >> 2
	longHeader.PacketNumberLength = b&0x03 + 1
}

func (longHeader *LongHeader) Write(packetBuffer *util.WriteBuffer) {
	b := longHeader.GetFirstByte()
	packetBuffer.WriteByte(b)
	packetBuffer.WriteBytes(longHeader.Version.ToBytes())
	packetBuffer.WriteUint8LengthBytes(longHeader.DestinationConnectionId)
	packetBuffer.WriteUint8LengthBytes(longHeader.SourceConnectionId)
	if longHeader.LongPacketType == packetTypeInitial {
		packetBuffer.WriteVariableLengthBytes(longHeader.Token)
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(longHeader.GetPacketNumberLength(), uint64(longHeader.PacketNumber))
			wr.WriteBytes(longHeader.PacketPayload)
		})
	}
	if longHeader.LongPacketType == packetTypeHandshake || longHeader.LongPacketType == packetTypeZeroRTT {
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(longHeader.GetPacketNumberLength(), uint64(longHeader.PacketNumber))
			wr.WriteBytes(longHeader.PacketPayload)
		})
	}
	if longHeader.LongPacketType == packetTypeRetry {
		packetBuffer.WriteBytes(longHeader.RetryToken)
		packetBuffer.WriteBytes(longHeader.RetryIntegrityTag)
	}
}
func (longHeader *LongHeader) Read(packetBuffer *util.ReadBuffer) error {
	readByte, err := packetBuffer.ReadByte()
	if err != nil {
		return err
	}
	longHeader.SetFirstByte(readByte)
	u32, err := packetBuffer.Read4U32()
	if err != nil {
		return err
	}
	longHeader.Version = util.VersionNumber(u32)
	_, data, err := packetBuffer.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	longHeader.DestinationConnectionId = data
	_, data, err = packetBuffer.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	longHeader.SourceConnectionId = data
	_, data, err = packetBuffer.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	longHeader.Token = data
	_, data, err = packetBuffer.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	longHeader.PacketPayload = data
	return nil
}


func NewLongHeader(longPacketType packetType, PacketPayload []byte, sendConfig *config.SendConfig) *LongHeader {
	var longHeader LongHeader
	longHeader.LongPacketType = longPacketType
	longHeader.IsLongHeader = true
	longHeader.Version = sendConfig.Version
	longHeader.DestinationConnectionId = sendConfig.ConnectionId
	longHeader.SourceConnectionId = []byte{}
	longHeader.Token = sendConfig.Token
	longHeader.PacketNumber = sendConfig.PacketNumber
	longHeader.PacketNumberLength = sendConfig.PacketNumber.GetPacketNumberLength()
	longHeader.PacketPayload = PacketPayload
	return &longHeader
}
