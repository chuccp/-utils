package wire

import "github.com/chuccp/utils/udp"

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

	var cryptoFrame CryptoFrame

	err := udp.UnPacketInitialPayload(packet.Header,&cryptoFrame)
	if err != nil {
		return err
	}

	return nil

}
