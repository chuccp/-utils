package udp

type packetType uint8

const (
	packetTypeInitial packetType = iota
	packetTypeZeroRTT
	packetTypeHandshake
	packetTypeRetry
	packetTypeVersionNegotiation
	packetTypeOneRTT
)


