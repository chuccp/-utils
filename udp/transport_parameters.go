package udp

import (
	"bytes"
	"time"
)
type transportParameterID byte
const (
	originalDestinationConnectionIDParameterID transportParameterID = 0x0
	maxIdleTimeoutParameterID                  transportParameterID = 0x1
	statelessResetTokenParameterID             transportParameterID = 0x2
	maxUDPPayloadSizeParameterID               transportParameterID = 0x3
	initialMaxDataParameterID                  transportParameterID = 0x4
	initialMaxStreamDataBidiLocalParameterID   transportParameterID = 0x5
	initialMaxStreamDataBidiRemoteParameterID  transportParameterID = 0x6
	initialMaxStreamDataUniParameterID         transportParameterID = 0x7
	initialMaxStreamsBidiParameterID           transportParameterID = 0x8
	initialMaxStreamsUniParameterID            transportParameterID = 0x9
	ackDelayExponentParameterID                transportParameterID = 0xa
	maxAckDelayParameterID                     transportParameterID = 0xb
	disableActiveMigrationParameterID          transportParameterID = 0xc
	preferredAddressParameterID                transportParameterID = 0xd
	activeConnectionIDLimitParameterID         transportParameterID = 0xe
	initialSourceConnectionIDParameterID       transportParameterID = 0xf
	retrySourceConnectionIDParameterID         transportParameterID = 0x10
	// https://datatracker.ietf.org/doc/draft-ietf-quic-datagram/
	maxDatagramFrameSizeParameterID transportParameterID = 0x20
	grease transportParameterID  = 0x3a
)
type TransportParameters struct {
	InitialMaxStreamDataBidiLocal  ByteCount
	InitialMaxStreamDataBidiRemote ByteCount
	InitialMaxStreamDataUni        ByteCount
	InitialMaxData                 ByteCount
	MaxBidiStreamNum               ByteCount
	MaxUniStreamNum                ByteCount
	MaxIdleTimeout                 time.Duration
	MaxUDPPayloadSize              ByteCount
	MaxAckDelay                    time.Duration
	DisableActiveMigration         bool
	ActiveConnectionIDLimit        ByteCount
	InitialSourceConnectionID      bool
	MaxDatagramFrameSize           ByteCount
	buffer *bytes.Buffer


}

func NewTransportParameters() *TransportParameters {
	transportParameters :=  &TransportParameters{buffer: new(bytes.Buffer)}
	transportParameters.InitialMaxStreamDataBidiLocal = 524288
	transportParameters.InitialMaxStreamDataBidiRemote = 524288
	transportParameters.InitialMaxStreamDataUni = 524288
	transportParameters.InitialMaxData = 786432
	transportParameters.MaxBidiStreamNum = 100
	transportParameters.MaxUniStreamNum = 100
	transportParameters.MaxIdleTimeout = time.Second*30
	transportParameters.MaxUDPPayloadSize = 1452
	transportParameters.DisableActiveMigration = false
	transportParameters.ActiveConnectionIDLimit = 4
	transportParameters.InitialSourceConnectionID = false
	transportParameters.MaxDatagramFrameSize = 0
	return transportParameters
}
func (p *TransportParameters)Bytes()[]byte{
	return p.buffer.Bytes()
}
func (p *TransportParameters)Init(){
	p.marshalByteParam(0x3a,0xcb)
	p.marshalVariableParam(initialMaxStreamDataBidiLocalParameterID, uint32(p.InitialMaxStreamDataBidiLocal))
	p.marshalVariableParam(initialMaxStreamDataBidiRemoteParameterID, uint32(p.InitialMaxStreamDataBidiRemote))
	p.marshalVariableParam(initialMaxStreamDataUniParameterID, uint32(p.InitialMaxStreamDataUni))
	p.marshalVariableParam(initialMaxDataParameterID, uint32(p.InitialMaxData))
	p.marshalVariableParam(initialMaxStreamsBidiParameterID, uint32(p.MaxBidiStreamNum))
	p.marshalVariableParam(initialMaxStreamsUniParameterID, uint32(p.MaxUniStreamNum))
	p.marshalVariableParam(maxIdleTimeoutParameterID, uint32(p.MaxIdleTimeout.Milliseconds()))
	p.marshalVariableParam(maxUDPPayloadSizeParameterID, uint32(p.MaxUDPPayloadSize))
	p.marshalByteParam(0x0b,0x1a)
	if p.DisableActiveMigration{
		p.marshalVariableParam(disableActiveMigrationParameterID,1)
	}else{
		p.marshalZeroParam(disableActiveMigrationParameterID)
	}
	p.marshalVariableParam(activeConnectionIDLimitParameterID, uint32(p.ActiveConnectionIDLimit))
	if p.InitialSourceConnectionID{
		p.marshalVariableParam(initialSourceConnectionIDParameterID,1)
	}else{
		p.marshalZeroParam(initialSourceConnectionIDParameterID)
	}
	p.marshalVariableParam(maxDatagramFrameSizeParameterID, uint32(p.MaxDatagramFrameSize))

}
func (p *TransportParameters) marshalVariableParam(id transportParameterID, val uint32) {
	p.buffer.WriteByte(byte(id))
	data:=VariableLengthBytes2(val)
	p.buffer.WriteByte(byte(len(data)))
	p.buffer.Write(data)
}

func (p *TransportParameters) marshalByteParam(id transportParameterID, val byte) {
	p.buffer.WriteByte(byte(id))
	p.buffer.WriteByte(1)
	p.buffer.WriteByte(val)
}
func (p *TransportParameters) marshalZeroParam(id transportParameterID) {
	p.buffer.WriteByte(byte(id))
	p.buffer.WriteByte(0)
}
