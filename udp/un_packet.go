package udp

import (
	"github.com/chuccp/utils/udp/tls"
	"github.com/chuccp/utils/udp/util"
	"github.com/chuccp/utils/udp/wire"
)

type Packet struct {
}

func UnPacket(data []byte, longHeader *LongHeader) error {
	fistByte := data[0]
	if (fistByte & 0x80) != 0 {
		rb := util.NewReadBuffer(data)
		err := longHeader.Read(rb)
		if err != nil {
			return err
		}
		if longHeader.LongPacketType == packetTypeInitial {
			err := UnPacketInitialPayload(longHeader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UnPacketInitialPayload(longHeader *LongHeader) error {
	rb := util.NewReadBuffer(longHeader.PacketPayload)
	u32, err := rb.ReadU8LengthU32(longHeader.PacketNumberLength)
	if err != nil {
		return err
	}
	longHeader.PacketNumber = util.PacketNumber(u32)
	for {
		readByte, err := rb.ReadByte()
		if err != nil {
			return err
		}
		if readByte == 0x6 {

			cryptoFrame,err := wire.ReadCryptoFrame(rb)
			if err != nil {
				return err
			}
			err = UnCryptoFramePayload(cryptoFrame)
			if err != nil {
				return err
			}
		}
		if rb.Buffered() == 0 {
			break
		}
	}

	return nil
}
func UnCryptoFramePayload(cryptoFrame *wire.CryptoFrame) error {
	rb := util.NewReadBuffer(cryptoFrame.Data)
	readByte, err := rb.ReadByte()
	if err != nil {
		return err
	}
	if tls.HandshakeType(readByte) == tls.ClientHelloType {
		_, err := tls.ReadClientHello(rb)
		if err != nil {
			return err
		}
		//log.Print(hello)
	}
	return err
}
