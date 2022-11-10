package config

import (
	"crypto/rand"
	"github.com/chuccp/utils/udp/util"
)

type ReceiveConfig struct {
}

type Config struct {
	PacketNumber                   util.PacketNumber
	Version                        util.VersionNumber
	ConnectionId                   []byte
	Token                          []byte
	TLSRandom                      []byte
	KeyExchanges                   []byte
	MaxUdpPayloadSize              uint32
	InitialMaxStreamDataBidiLocal  uint32
	InitialMaxStreamDataBidiRemote uint32
	InitialMaxData                 uint32
	MaxIdleTimeout                 uint32
	MaxDatagramFrameSize           uint32
	InitialMaxStreamsBidi          uint32
	InitialMaxStreamsUni           uint32
}

func NewReceiveConfig() *Config {
	return &Config{}
}

func CreateSendConfig() *Config {
	TLSRandom := make([]byte, 32)
	rand.Read(TLSRandom)
	KeyExchanges := make([]byte, 32)
	rand.Read(KeyExchanges)
	ConnectionId := make([]byte, 16)
	rand.Read(ConnectionId)
	token := make([]byte, 16)
	rand.Read(token)
	return &Config{PacketNumber: 0, Version: util.Version1, ConnectionId: ConnectionId, Token: token, TLSRandom: TLSRandom, KeyExchanges: KeyExchanges}
}
