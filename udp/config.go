package udp




type SendConfig struct {
	PacketNumber PacketNumber
	Version VersionNumber
	ConnectionId []byte
	Token []byte
}

func (send *SendConfig) NewSendConfig()*SendConfig  {

	return &SendConfig{PacketNumber:0,}
}