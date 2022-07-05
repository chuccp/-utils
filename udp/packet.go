package udp

type packetType uint8

const (
	packetTypeInitial packetType = iota
	packetTypeZeroRTT
	packetTypeHandshake
	packetTypeRetry
	packetTypeVersionNegotiation
	packetTypeOneRTT
)

func Packet(longHeader *LongHeader)  {
	var packetBuffer = NewPacketBuffer()
	b := longHeader.GetFirstByte()
	packetBuffer.WriteByte(b)
	packetBuffer.WriteBytes(longHeader.Version.ToBytes())
	packetBuffer.WriteByte(longHeader.ConnectionIdLength)
	packetBuffer.WriteBytes(longHeader.ConnectionId)
	if longHeader.LongPacketType==packetTypeInitial{
		packetBuffer.WriteVariableLength(longHeader.TokenVariableLength)
		packetBuffer.WriteBytes(longHeader.Token)
		packetBuffer.WriteVariableLength(longHeader.LengthVariable)
		packetBuffer.WriteUint64(longHeader.GetPacketNumberLength(), uint64(longHeader.PacketNumber))
		packetBuffer.WriteBytes(longHeader.PacketPayload)
	}
	if longHeader.LongPacketType==packetTypeHandshake || longHeader.LongPacketType==packetTypeZeroRTT{
		packetBuffer.WriteVariableLength(longHeader.LengthVariable)
		packetBuffer.WriteUint64(longHeader.GetPacketNumberLength(), uint64(longHeader.PacketNumber))
		packetBuffer.WriteBytes(longHeader.PacketPayload)
	}
	if longHeader.LongPacketType==packetTypeRetry{
		packetBuffer.WriteBytes(longHeader.RetryToken)
		packetBuffer.WriteBytes(longHeader.RetryIntegrityTag)
	}
}
