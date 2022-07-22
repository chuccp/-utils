package wire

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
	err = UnPacketInitialPayload(packet.Header,&cryptoFrame)
	if err != nil {
		return err
	}
	return nil

}
