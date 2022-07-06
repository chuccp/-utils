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

func Packet(longHeader *LongHeader,packetBuffer *PacketWriteBuffer)  {
	b := longHeader.GetFirstByte()
	packetBuffer.WriteByte(b)
	packetBuffer.WriteBytes(longHeader.Version.ToBytes())
	packetBuffer.WriteByte(longHeader.DestinationConnectionIdLength)
	packetBuffer.WriteBytes(longHeader.DestinationConnectionId)
	packetBuffer.WriteByte(longHeader.SourceConnectionIdLength)
	packetBuffer.WriteBytes(longHeader.SourceConnectionId)
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
