package udp

import "github.com/chuccp/utils/math"

func UnPacket(data []byte) {
	fistByte := data[0]
	if (fistByte & 0x80) !=0 {
		var longHeader LongHeader
		UnPacketLongHeader(data, &longHeader)
	}

}
func UnPacketLongHeader(data []byte, longHeader *LongHeader) error {
	packetReadBuffer := NewPacketReadBuffer(data)

	readByte, err := packetReadBuffer.ReadByte()
	if err != nil {
		return err
	}
	longHeader.SetFirstByte(readByte)
	u32, err := packetReadBuffer.Read4U32()
	if err != nil {
		return err
	}
	longHeader.Version = VersionNumber(u32)

	u8, bytes, err := packetReadBuffer.ReadU8Bytes()
	if err != nil {
		return err
	}
	longHeader.DestinationConnectionIdLength = u8
	longHeader.DestinationConnectionId = bytes
	u8, bytes, err = packetReadBuffer.ReadU8Bytes()
	if err != nil {
		return err
	}
	longHeader.SourceConnectionIdLength = u8
	longHeader.SourceConnectionId = bytes

	if longHeader.LongPacketType == packetTypeInitial {
		longHeader.TokenVariableLength, longHeader.Token, err = packetReadBuffer.ReadVariableLengthBytes()
		if err != nil {
			return err
		}
		longHeader.LengthVariable, longHeader.PacketPayload, err = packetReadBuffer.ReadVariableLengthBytes()
		if err != nil {
			return err
		}
		pn := int(longHeader.PacketNumberLength+1)
		data := longHeader.PacketPayload[0:pn]
		longHeader.PacketNumber = PacketNumber(math.U32BE0To4(data, uint8(pn)))
		longHeader.PacketPayload = longHeader.PacketPayload[pn:]
		return err
	}
	if longHeader.LongPacketType == packetTypeHandshake || longHeader.LongPacketType == packetTypeZeroRTT {

		longHeader.LengthVariable, longHeader.PacketPayload, err = packetReadBuffer.ReadVariableLengthBytes()
		if err != nil {
			return err
		}
		pn := int(longHeader.PacketNumberLength)
		data := longHeader.PacketPayload[0:pn]
		longHeader.PacketNumber = PacketNumber(math.U32BE0To4(data, uint8(pn)))
		longHeader.PacketPayload = longHeader.PacketPayload[int(longHeader.PacketNumberLength):]
		return err
	}
	if longHeader.LongPacketType == packetTypeRetry {
		num := packetReadBuffer.Offset()
		nLen := len(data)
		mLen := uint16(nLen) - 16 - num
		longHeader.RetryToken, err = packetReadBuffer.ReadU32Bytes(uint32(mLen))
		if err != nil {
			return err
		}
		longHeader.RetryIntegrityTag, err = packetReadBuffer.ReadBytes(16)
		if err != nil {
			return err
		}
	}

	return nil

}
