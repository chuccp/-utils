package tls

import (
	"github.com/chuccp/utils/udp/util"
)

const (
	MaxUdpPayloadSizeType byte = 0x03
	InitialMaxDataType byte = 0x04
	InitialMaxStreamDataBidiLocalType byte = 0x05
	InitialMaxStreamDataBidiRemoteType byte = 0x06
	InitialMaxStreamsBidi byte  = 0x08
	InitialMaxStreamsUni byte  = 0x09
	MaxIdleTimeout byte  = 0x01

)

type TransportParameters struct {
	TransportParameterMap map[byte]*TransportParameter
}

func NewTransportParameters() *TransportParameters {
	transportParameters:=make(map[byte]*TransportParameter)
	return &TransportParameters{TransportParameterMap:transportParameters}
}
func (tps *TransportParameters)SetValue(Type byte,value uint32){
	TransportParameter:=NewTransportParameter(Type,value )
	tps.TransportParameterMap[TransportParameter.Type]  = TransportParameter
}
func (tps *TransportParameters) Bytes(write *util.WriteBuffer)  {
	for _, parameter := range tps.TransportParameterMap {
		parameter.Bytes(write)
	}
}


type TransportParameter struct {
	Type byte
	Length uint8
	Value []byte
}

func NewTransportParameter(Type byte,value uint32) *TransportParameter {
	data:=util.VariableLengthToBytes(value)
	return &TransportParameter{Type:Type,Length: uint8(len(data)),Value: data}
}
func (tp *TransportParameter) Bytes(write *util.WriteBuffer)  {
	write.WriteByte(tp.Type)
	write.WriteByte(tp.Length)
	if tp.Length>0{
		write.WriteBytes(tp.Value)
	}
}