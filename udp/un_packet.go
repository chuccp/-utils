package udp

import (
	"github.com/chuccp/utils/udp/util"
	"log"
)

func UnPacket(data []byte) error {
	fistByte := data[0]
	if (fistByte & 0x80) != 0 {
		var longHeader LongHeader
		return UnPacketLongHeader(data, &longHeader)
	}
	return nil
}
func UnPacketLongHeader(data []byte, longHeader *LongHeader) error {
	rb := util.NewReadBuffer(data)
	err := rb.ReadByteBuff(func(b byte, buffer *util.ReadBuffer) error {
		longHeader.SetFirstByte(b)
		buffer.ReadUInt32Buff(func(u uint32, buffer *util.ReadBuffer) error {
			longHeader.Version = util.VersionNumber(u)
			buffer.ReadUint8LengthBytesBuff(func(bytes []byte, buffer *util.ReadBuffer) error {
				longHeader.DestinationConnectionId = bytes
				buffer.ReadUint8LengthBytesBuff(func(bytes []byte, buffer *util.ReadBuffer) error {
					longHeader.SourceConnectionId = bytes
					buffer.ReadVariableLengthBytesBuff(func(bytes []byte, buffer *util.ReadBuffer) error {
						longHeader.Token = bytes
						_, data, err := buffer.ReadVariableLengthBytes()
						if err != nil {
							return err
						} else {
							log.Print("===========",data)

							

						}
						return nil
					})
					return nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}
