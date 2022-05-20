package wire2

type PacketType uint8
const (
	// PacketTypeInitial is the packet type of an Initial packet
	PacketTypeInitial PacketType =  iota
	// PacketTypeRetry is the packet type of a Retry packet
	PacketTypeRetry
	// PacketTypeHandshake is the packet type of a Handshake packet
	PacketTypeHandshake
	// PacketType0RTT is the packet type of a 0-RTT packet
	PacketType0RTT
)
type HandshakeType uint8

const (
	ClientHello HandshakeType = 1
	ServerHello  HandshakeType = 2
)

type ExtensionType uint16

const (
	StatusRequest ExtensionType = 5
)

type ParameterType uint8

type LongHeader struct {
	FixedBit bool
	PacketType PacketType
	Reserved byte
	PacketNumberLength byte
	Version []byte
	DestinationConnectionIdLength byte
	DestinationConnectionId []byte
	SourceConnectionIdLength byte
	SourceConnectionId []byte
	TokenLength byte
	Token []byte
	Length uint32
	PackageNum uint32
}



type ClientInitial struct {
	HandshakeType byte
	Length uint32
	Version []byte
	Random []byte
	SessionIdLength uint8
	SessionId []byte
	CipherSuiteLength uint16
	CipherSuites []byte
	CompressMethodsLength byte
	CompressMethods []byte
}
type Extension struct {
	ExtensionType ExtensionType
	Length uint16
	Value []byte
}
type Extensions struct {
	Length uint16
}

type Parameter struct {
	ParameterType ParameterType
	Length uint8
	Value []byte
}
type TransportParameters struct {
	Parameters []Parameter
}