package tls

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

const (
	MaxUdpPayloadSizeType              byte = 0x03
	InitialMaxDataType                 byte = 0x04
	InitialMaxStreamDataBidiLocalType  byte = 0x05
	InitialMaxStreamDataBidiRemoteType byte = 0x06
	InitialMaxStreamsBidi              byte = 0x08
	InitialMaxStreamsUni               byte = 0x09
	MaxIdleTimeout                     byte = 0x01
)

type TransportParameters struct {
	TransportParameterMap map[byte]*TransportParameter
}

func NewTransportParameters(sendConfig *config.SendConfig) *TransportParameters {
	transportParameterMap := make(map[byte]*TransportParameter)
	transportParameters := &TransportParameters{TransportParameterMap: transportParameterMap}
	transportParameters.SetValue(MaxUdpPayloadSizeType, sendConfig.MaxUdpPayloadSize)
	transportParameters.SetValue(InitialMaxStreamDataBidiRemoteType, sendConfig.InitialMaxStreamDataBidiRemote)
	transportParameters.SetValue(InitialMaxStreamDataBidiLocalType, sendConfig.InitialMaxStreamDataBidiLocal)
	transportParameters.SetValue(InitialMaxDataType, sendConfig.InitialMaxData)
	transportParameters.SetValue(MaxIdleTimeout, sendConfig.MaxIdleTimeout)
	transportParameters.SetValue(InitialMaxStreamsBidi, sendConfig.InitialMaxStreamsBidi)
	transportParameters.SetValue(InitialMaxStreamsUni, sendConfig.InitialMaxStreamsUni)
	return transportParameters
}
func (tps *TransportParameters) SetValue(Type byte, value uint32) {
	TransportParameter := NewTransportParameter(Type, value)
	tps.TransportParameterMap[TransportParameter.Type] = TransportParameter
}
func (tps *TransportParameters) Write(write *util.WriteBuffer) {
	for _, parameter := range tps.TransportParameterMap {
		parameter.Write(write)
	}
}

type TransportParameter struct {
	Type   byte
	Length uint8
	Value  []byte
}

func NewTransportParameter(Type byte, value uint32) *TransportParameter {
	data := util.VariableLengthToBytes(value)
	return &TransportParameter{Type: Type, Length: uint8(len(data)), Value: data}
}
func (tp *TransportParameter) Write(write *util.WriteBuffer) {
	write.WriteByte(tp.Type)
	write.WriteByte(tp.Length)
	if tp.Length > 0 {
		write.WriteBytes(tp.Value)
	}
}
