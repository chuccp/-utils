package tls

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

const (
	KeyShareType uint16 = 51
)

type Extensions struct {
	Extensions []*Extension
}

func (es *Extensions) SetKeyShare(KeyExchanges []byte) {
	ks := NewKeyShare(KeyExchanges)
	ex := NewExtension(KeyShareType, ks)
	es.addExtensions(ex)
}

func (es *Extensions) SetTransportParameters(sendConfig *config.SendConfig) {
	ntp:=NewTransportParameters()
	ntp.SetValue(MaxUdpPayloadSizeType,sendConfig.MaxUdpPayloadSize)
	ntp.SetValue(InitialMaxStreamDataBidiRemoteType,sendConfig.InitialMaxStreamDataBidiRemote)
	ntp.SetValue(InitialMaxStreamDataBidiLocalType,sendConfig.InitialMaxStreamDataBidiLocal)
	ntp.SetValue(InitialMaxDataType,sendConfig.InitialMaxData)
	ntp.SetValue(MaxIdleTimeout,sendConfig.MaxIdleTimeout)
	ntp.SetValue(InitialMaxStreamsBidi,sendConfig.InitialMaxStreamsBidi)
	ntp.SetValue(InitialMaxStreamsUni,sendConfig.InitialMaxStreamsUni)
	ex := NewExtension(KeyShareType, ntp)
	es.addExtensions(ex)
}


func (es *Extensions) addExtensions(ex *Extension) {
	es.Extensions = append(es.Extensions, ex)
}

func NewExtensions() *Extensions {
	return &Extensions{Extensions: make([]*Extension, 0)}
}

func (es *Extensions) Write(write *util.WriteBuffer) {
	for _, extension := range es.Extensions {
		write.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
			extension.Write(wr)
		})
	}
}

type Extension struct {
	Type   uint16
	Buffer util.BufferWrite
}

func NewExtension(Type uint16, Buffer util.BufferWrite) *Extension {
	return &Extension{Type: Type, Buffer: Buffer}
}
func (e *Extension) Write(write *util.WriteBuffer) {
	write.WriteUint16(e.Type)
	write.WriteUint16LengthBuff(func(write *util.WriteBuffer) {
		e.Buffer.Write(write)
	})
}
