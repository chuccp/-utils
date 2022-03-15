package udp

type Package struct {
	headerFrom         bool
	fixedBig           bool
	longPacketType     byte
	reservedBits       byte
	packetNumberLength byte
	version            uint32
	desConnIdLength    uint8
	desConnId          []byte
	sourceConnLength   uint8
	sourceConn         []byte
	tokenLength        uint8
	token              []byte
	entryContextLength uint8
	packetNumber uint8
	entryContext []byte
}
