package wire

import (
	"github.com/chuccp/utils/udp/tls"
	"github.com/chuccp/utils/udp/util"
)

type HandshakeStatusType uint8

const (
	WaitInitialHandshake HandshakeStatusType = iota
	FinishHandshake
)

type ServerHandshake struct {
	HandshakeStatus HandshakeStatusType
}

func NewServerHandshake() *ServerHandshake {
	return &ServerHandshake{}
}
func (serverHandshake *ServerHandshake) Handle(packet *ReceivePacket) error {
	err := packet.Header.ParseExLongHeader(packet.Data)
	if err != nil {
		return err
	}
	var cryptoFrame CryptoFrame
	err = UnPacketCryptoFrame(packet.Header.PacketPayload,&cryptoFrame)
	if err != nil {
		return err
	}
	if cryptoFrame.Data[0]==0x01{

		var ch  tls.ClientHello
		err := tls.UnClientHelloHandshake(cryptoFrame.Data, &ch)
		if err != nil {
			return err
		}

		return nil
	}
	return util.FormatError

}
