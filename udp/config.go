package udp

type SendConfig struct {
	PacketNumber PacketNumber
	Version VersionNumber
	ConnectionId []byte
	Token []byte
}

func NewSendConfig(ConnectionId []byte)*SendConfig  {

	return &SendConfig{PacketNumber:0,Version:Version1,ConnectionId:ConnectionId,Token:[]byte{}}
}