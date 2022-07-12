package config

import (
	"crypto/rand"
	"github.com/chuccp/utils/udp/util"
)

type SendConfig struct {
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
	InitialMaxStreamsBidi uint32
	InitialMaxStreamsUni uint32
}

func NewSendConfig(ConnectionId []byte)*SendConfig {
	TLSRandom:=make([]byte,32)
	rand.Read(TLSRandom)
	return &SendConfig{PacketNumber: 0,Version: util.Version1,ConnectionId:ConnectionId,Token:[]byte{},TLSRandom:TLSRandom}
}