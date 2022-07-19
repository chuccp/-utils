package wire

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
	IsLongHeader bool
	LongPacketType    PacketType

}
