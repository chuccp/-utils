package wire

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/udp/util"
	"math"
)

type InitialParameter struct {
	ConnectionID []byte
	Token        []byte
	PacketNum    ByteCount
	Random       []byte
	PacketType   PacketType

}

type LongHeaderPackage struct {
	headerFrom         bool
	fixedBig           bool
	longPacketType     byte
	reservedBits       byte
	packetNumberLength uint8
	version            []byte
	desConnIdLength    uint8
	desConnId          []byte
	sourceConnLength   uint8
	sourceConnId       []byte
	tokenLength        []byte
	token              []byte
	length             []byte
	packetNumber       []byte
	payload            []byte
}

func (longHeaderPackage *LongHeaderPackage) Bytes()[]byte  {
	var buff = new(bytes.Buffer)
	aead:=NewInitialAEAD(longHeaderPackage.desConnId)
	var headType byte= 0
	if longHeaderPackage.headerFrom{
		headType = headType|128
	}
	if longHeaderPackage.fixedBig{
		headType = headType|64
	}
	headType = headType|(longHeaderPackage.longPacketType<<4)
	headType = headType|(longHeaderPackage.reservedBits<<2)
	headType = headType|(longHeaderPackage.packetNumberLength-1)
	buff.WriteByte(headType)
	buff.Write(longHeaderPackage.version)
	buff.WriteByte(longHeaderPackage.desConnIdLength)
	buff.Write(longHeaderPackage.desConnId)
	buff.WriteByte(longHeaderPackage.sourceConnLength)
	buff.Write(longHeaderPackage.sourceConnId)
	buff.Write(longHeaderPackage.tokenLength)
	buff.Write(longHeaderPackage.token)

	var payloadLen = 1227-16

	buff.Write(util.VariableLengthBytes(1227))


	rdata:=buff.Bytes()


	buff.Write(longHeaderPackage.packetNumber)
	exLen:=buff.Len()
	additionalData:=buff.Bytes()

	log.Info(ConnectionID(additionalData))

	buff.Write(longHeaderPackage.payload)
	data := buff.Bytes()
	text:=data[exLen:]
	nonceBuf:=make([]byte, aead.aead.NonceSize())


	copy(nonceBuf[len(nonceBuf)-int(longHeaderPackage.packetNumberLength):],longHeaderPackage.packetNumber)
	for i, b := range nonceBuf[len(nonceBuf)-8:] {
		aead.iv[4+i] ^= b
	}


	readLen:=payloadLen-int(longHeaderPackage.packetNumberLength)

	payload:=make([]byte,readLen)

	copy(payload[readLen-len(text):],text)


	ciphertext:=aead.aead.Seal([]byte{},aead.iv,payload,additionalData)


	sample:=ciphertext[4-int(longHeaderPackage.packetNumberLength):20-int(longHeaderPackage.packetNumberLength)]
	mask:=make([]byte, aead.block.BlockSize())

	aead.block.Encrypt(mask,sample)

	rdata[0] ^= mask[0] & 0xf

	for i := range longHeaderPackage.packetNumber {
		longHeaderPackage.packetNumber[i] ^= mask[i+1]
	}
	var buff2 = new(bytes.Buffer)
	buff2.Write(rdata)
	buff2.Write(longHeaderPackage.packetNumber)
	buff2.Write(ciphertext)
	return buff2.Bytes()
}






func Initial(initialParameter *InitialParameter)*LongHeaderPackage  {
	return LongHeaderPacket(initialParameter)
}

func LongHeaderPacket(initialParameter *InitialParameter)*LongHeaderPackage  {
	var packetNumberLength uint8 = 2
	if initialParameter.PacketNum >math.MaxUint16{
		packetNumberLength = 4
	}
	pn:=make([]byte,8)
	binary.BigEndian.PutUint64(pn, uint64(initialParameter.PacketNum))
	pn = pn[8-packetNumberLength:]
	ch:= NewClientHello(initialParameter.Random)
	dl:=ch.Bytes()
	crypto:=NewCrypto(0, dl)
	tokenLen := len(initialParameter.Token)
	return &LongHeaderPackage{
		headerFrom:true,fixedBig:true,longPacketType: byte(initialParameter.PacketType),reservedBits:0,
		packetNumberLength:packetNumberLength,version:[]byte{0,0,0,1},
		desConnIdLength: uint8(len(initialParameter.ConnectionID)),desConnId:initialParameter.ConnectionID,
		sourceConnLength:0, sourceConnId:[]byte{},tokenLength:util.VariableLengthBytes(uint32(tokenLen)),token:initialParameter.Token,packetNumber:pn,
		payload:crypto.Bytes(),
	}
}