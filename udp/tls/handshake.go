package tls

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

type HandshakeType byte

const (
	ClientHelloType HandshakeType = 1
	ServerHelloType HandshakeType = 2
)

type ClientHello struct {
	HandshakeType HandshakeType
	Length        uint16
	Version       []byte
	Random        []byte
	//SessionIdLength          uint8
	SessionId []byte
	//CipherSuitesLength       uint8
	CipherSuites             *CipherSuites
	CompressionMethodsLength uint8
	CompressionMethods       []byte
	Extensions               *Extensions
}

func NewClientHello(sendConfig *config.SendConfig) *ClientHello {

	Extensions := NewExtensions()
	Extensions.SetKeyShare(sendConfig.KeyExchanges)
	Extensions.SetTransportParameters(sendConfig)

	return &ClientHello{HandshakeType: 1, Version: []byte{03, 07},
		Random: sendConfig.TLSRandom, SessionId: []byte{}, CompressionMethods: []byte{0},
		CipherSuites: NewCipherSuites(), Extensions: Extensions}
}

func (ch *ClientHello) Write(write *util.WriteBuffer) {
	write.WriteByte(byte(ch.HandshakeType))
	write.WriteUint24LengthBuff(func(w *util.WriteBuffer) {
		w.WriteBytes(ch.Version)
		w.WriteBytes(ch.Random[0:32])
		w.WriteUint8LengthBytes(ch.SessionId)
		w.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
			ch.CipherSuites.Write(wr)
		})
		w.WriteUint8LengthBytes(ch.CompressionMethods)



		w.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			ch.Extensions.Write(wr)
		})
	})
}
func (ch *ClientHello) Read(read *util.ReadBuffer) error {

	length, err := read.ReadUint24Length()
	if err != nil {
		return err
	}
	data, err := read.ReadU32Bytes(length)
	if err != nil {
		return err
	}
	rd := util.NewReadBuffer(data)
	ch.Version, err = rd.ReadBytes(2)
	if err != nil {
		return err
	}
	ch.Random, err = rd.ReadBytes(32)
	if err != nil {
		return err
	}
	_, ch.SessionId, err = rd.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	_, data, err = rd.ReadU16LengthBytes()
	if err != nil {
		return err
	}
	cs := util.NewReadBuffer(data)
	ch.CipherSuites, err = ReadCipherSuites(cs)
	if err != nil {
		return err
	}
	_, ch.CompressionMethods, err = rd.ReadU8LengthBytes()
	if err != nil {
		return err
	}
	_, exData, err := rd.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	ex := util.NewReadBuffer(exData)
	ch.Extensions, err = ReadExtensions(ex)
	if err != nil {
		return err
	}
	return nil
}

func UnPacketClientHelloHandshake(data []byte,ch  *ClientHello)error{
	rd:=util.NewReadBuffer(data[1:])
	return ch.Read(rd)
}
