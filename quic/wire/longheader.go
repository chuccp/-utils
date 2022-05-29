package wire

import (
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/math"
	"github.com/chuccp/utils/quic/util"
	"log"
)

type PacketType uint8

const (
	// PacketTypeInitial is the packet type of an Initial packet
	PacketTypeInitial PacketType = iota
	// PacketTypeRetry is the packet type of a Retry packet
	PacketTypeRetry
	// PacketTypeHandshake is the packet type of a Handshake packet
	PacketTypeHandshake
	// PacketType0RTT is the packet type of a 0-RTT packet
	PacketType0RTT
)

type HandshakeType uint8

const (
	ClientHelloType HandshakeType = 1
	ServerHelloType HandshakeType = 2
)



type ParameterType uint8

type LongPackage struct {
	firstByte                     byte
	IsLongHeader                  bool
	FixedBit                      bool
	PacketType                    PacketType
	Reserved                      byte
	PacketNumberLength            byte
	Version                       uint32
	DestinationConnectionIdLength byte
	DestinationConnectionId       []byte
	SourceConnectionIdLength      byte
	SourceConnectionId            []byte
	TokenLength                   byte
	Token                         []byte
	Length                        uint32
	PackageNum                    uint32
	aead                          *AEAD
	PlayLoad                      []byte
}





type Parameter struct {
	ParameterType ParameterType
	Length        uint8
	Value         []byte
}
type TransportParameters struct {
	Parameters []Parameter
}
func ParseInitialHeader(data []byte) (longPackage *LongPackage, err error) {
	stream := io.NewReadBytesStream(data)
	longPackage = &LongPackage{}
	longPackage.firstByte, err = stream.ReadByte()
	log.Println("===1", stream.Offset())
	if err != nil {
		return nil, err
	}
	longPackage.Version, err = stream.Read4Uint32()
	log.Println("===2", stream.Offset())
	if err != nil {
		return nil, err
	}
	readValue := util.NewReadValue(stream)
	longPackage.DestinationConnectionId, longPackage.DestinationConnectionIdLength, err = readValue.ReadUint8()
	log.Println("===3", stream.Offset())
	if err != nil {
		return nil, err
	}
	longPackage.SourceConnectionId, longPackage.SourceConnectionIdLength, err = readValue.ReadUint8()
	log.Println("===4", stream.Offset())
	if err != nil {
		return nil, err
	}
	longPackage.Token, longPackage.TokenLength, err = readValue.ReadUint8()
	log.Println("===5", stream.Offset())
	if err != nil {
		return nil, err
	}
	longPackage.Length, err = readValue.ReadVariableValueLength()

	if err != nil {
		return nil, err
	}
	/**
	解析packageNum
	*/
	offset := stream.Offset()
	longPackage.aead = NewInitialAEAD(longPackage.DestinationConnectionId)
	longPackage.PacketNumberLength, longPackage.PackageNum, err = longPackage.ParsePackageNum(data, offset, uint16(longPackage.Length))
	if err != nil {
		return nil, err
	}
	exLen := offset + uint16(longPackage.PacketNumberLength)
	additionalData := data[:exLen]
	ciphertext := data[exLen : offset+uint16(longPackage.Length)]
	data, err = longPackage.aead.aead.Open([]byte{}, longPackage.aead.iv, ciphertext, additionalData)
	if err == nil {
		longPackage.PlayLoad = data
		return longPackage, nil
	}
	return nil, err
}
func ParseInitial(data []byte) (longPackage *LongPackage, err error) {
	longPackage, err = ParseInitialHeader(data)
	frame, err := ParseFrame(longPackage.PlayLoad)
	if err != nil {
		return nil, err
	}

	_, err = ParseHandshake(frame)
	if err != nil {
		return nil, err
	}


	return longPackage, err
}
func (longPackage *LongPackage) ParsePackageNum(data []byte, offset uint16, length uint16) (byte, uint32, error) {
	payload := data[offset : offset+length]
	origPNBytes := make([]byte, 4)
	copy(origPNBytes, payload[0:4])
	sample := payload[4:20]
	param2 := &longPackage.firstByte
	param3 := payload[0:4]
	mask := make([]byte, longPackage.aead.block.BlockSize())
	longPackage.aead.block.Encrypt(mask, sample)
	*param2 ^= mask[0] & 0xf
	for i := range param3 {
		param3[i] ^= mask[i+1]
	}
	pageNumLength := longPackage.firstByte&0x3 + 1
	copy(payload[pageNumLength:4], origPNBytes[pageNumLength:4])
	u := param3[0:pageNumLength]
	return pageNumLength, math.U32BE0To4(u, uint8(pageNumLength)), nil
}

func ParseLongPackage(data []byte) (*LongPackage, error) {
	packetType := PacketType(data[0] & 30 >> 4)
	if packetType == PacketTypeInitial {
		initial, err := ParseInitial(data)
		if err != nil {
			return initial, err
		}
	}
	return nil, util.ProtocolError
}
