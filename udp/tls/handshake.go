package tls

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

type ClientHello struct {
	HandshakeType            byte
	Length                   uint16
	Version                  []byte
	Random                   []byte
	//SessionIdLength          uint8
	SessionId                []byte
	//CipherSuitesLength       uint8
	CipherSuites             *CipherSuites
	CompressionMethodsLength uint8
	CompressionMethods       []byte
	Extensions               *Extensions
}

func NewClientHello(sendConfig *config.SendConfig) *ClientHello {

	Extensions:=NewExtensions()
	Extensions.SetKeyShare(sendConfig.KeyExchanges)
	Extensions.SetTransportParameters(sendConfig)

	return &ClientHello{HandshakeType:1,Version:[]byte{03,07},
		Random: sendConfig.TLSRandom,SessionId:[]byte{},CompressionMethods:[]byte{0},
		CipherSuites:NewCipherSuites(),Extensions:Extensions}
}


func (ch *ClientHello) Bytes(write *util.WriteBuffer) {
	write.WriteByte(ch.HandshakeType)
	write.WriteUint24LengthBuff(func(w *util.WriteBuffer) {
		w.WriteBytes(ch.Version)
		w.WriteBytes(ch.Random[0:32])
		w.WriteUint8LengthBytes(ch.SessionId)
		w.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
			ch.CipherSuites.Bytes(wr)
		})
		w.WriteUint8LengthBytes(ch.CompressionMethods)
		w.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
			ch.Extensions.Bytes(wr)
		})
	})
}
