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
	DestinationConnectionId       []byte

	//SourceConnectionIdLength uint8
	SourceConnectionId       []byte

	//TokenVariableLength uint32
	Token               []byte
	//LengthVariable      uint32
	PacketNumber  util.PacketNumber
	PacketPayload util.Buffer

	RetryToken        []byte
	RetryIntegrityTag []byte
}

func (h *LongHeader) GetPacketNumberLength() uint8 {
	return h.PacketNumberLength + 1
}
func (h *LongHeader) GetFirstByte() byte {
	var b byte = 0
	if h.IsLongHeader {
		b = b | 0x80
	}
	if h.FixedBit {
		b = b | 0x40
	}
	b = b | (uint8(h.LongPacketType) << 4)
	b = b | (h.ReservedBits << 2)
	b = b | (h.PacketNumberLength)
	return b
}
func (h *LongHeader) SetFirstByte(b byte) {
	h.IsLongHeader = b&0x80 > 0
	h.FixedBit = b&0x40 > 0
	h.LongPacketType = packetType(b&0x30>>4)
	h.ReservedBits = 0x0c>>2
	h.PacketNumberLength = b&0x03
}

func (h *LongHeader) Bytes(packetBuffer *util.WriteBuffer)  {
	b := h.GetFirstByte()
	packetBuffer.WriteByte(b)
	packetBuffer.WriteBytes(h.Version.ToBytes())
	packetBuffer.WriteUint8LengthBytes(h.DestinationConnectionId)
	packetBuffer.WriteUint8LengthBytes(h.SourceConnectionId)
	if h.LongPacketType==packetTypeInitial{
		packetBuffer.WriteVariableLengthBytes(h.Token)
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(h.GetPacketNumberLength(), uint64(h.PacketNumber))
			h.PacketPayload.Bytes(wr)
		})
	}
	if h.LongPacketType==packetTypeHandshake || h.LongPacketType==packetTypeZeroRTT{
		packetBuffer.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteLenUint64(h.GetPacketNumberLength(), uint64(h.PacketNumber))
			h.PacketPayload.Bytes(wr)
		})
	}
	if h.LongPacketType==packetTypeRetry{
		packetBuffer.WriteBytes(h.RetryToken)
		packetBuffer.WriteBytes(h.RetryIntegrityTag)
	}
}
func NewLongHeader(longPacketType packetType, PlayLoad util.Buffer, sendConfig *config.SendConfig) *LongHeader {
	var longHeader LongHeader
	longHeader.LongPacketType = longPacketType
	longHeader.IsLongHeader = true
	longHeader.Version = sendConfig.Version
	longHeader.DestinationConnectionId = sendConfig.ConnectionId
	longHeader.SourceConnectionId = []byte{}
	longHeader.Token = sendConfig.Token
	longHeader.PacketNumber = sendConfig.PacketNumber
	longHeader.PacketPayload = PlayLoad
	return &longHeader
}
