package udp

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/log"
	"math/rand"
)

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
	headType = headType|longHeaderPackage.packetNumberLength
	buff.WriteByte(headType)
	buff.Write(longHeaderPackage.version)
	buff.WriteByte(longHeaderPackage.desConnIdLength)
	buff.Write(longHeaderPackage.desConnId)
	buff.WriteByte(longHeaderPackage.sourceConnLength)
	buff.Write(longHeaderPackage.sourceConnId)
	buff.Write(longHeaderPackage.tokenLength)
	buff.Write(longHeaderPackage.token)
	buff.Write(VariableLengthBytes(uint32(len(longHeaderPackage.payload))+ 16))

	rdata:=buff.Bytes()

	buff.Write(longHeaderPackage.packetNumber)
	exLen:=buff.Len()
	additionalData:=buff.Bytes()
	buff.Write(longHeaderPackage.payload)
	data := buff.Bytes()
	text:=data[exLen:]
	nonceBuf:=make([]byte, aead.aead.NonceSize())
	copy(nonceBuf[len(nonceBuf)-int(longHeaderPackage.packetNumberLength):],longHeaderPackage.packetNumber)
	for i, b := range nonceBuf[len(nonceBuf)-8:] {
		aead.iv[4+i] ^= b
	}
	log.Info("!!!", len(text))
	ciphertext:=aead.aead.Seal([]byte{},aead.iv,text,additionalData)
	log.Info("!!!", len(ciphertext))
	sample:=ciphertext[4-int(longHeaderPackage.packetNumberLength):20-int(longHeaderPackage.packetNumberLength)]

	mask:=make([]byte, aead.block.BlockSize())
	aead.block.Decrypt(mask,sample)

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


func GenerateConnectionID(len int) (ConnectionID, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return ConnectionID(b), nil
}

func Initial(DesConnectionID []byte)*LongHeaderPackage  {
	return LongHeaderPacket(DesConnectionID,PacketTypeInitial)
}

func LongHeaderPacket(DesConnectionID []byte,packetType PacketType)*LongHeaderPackage  {

	var packetNumberLength uint8 = 1
	var packetNumber uint64=1
	pn:=make([]byte,8)
	binary.BigEndian.PutUint64(pn,packetNumber)
	pn = pn[8-packetNumber:]
	return &LongHeaderPackage{
		headerFrom:true,fixedBig:true,longPacketType: byte(packetType),reservedBits:0,
		packetNumberLength:packetNumberLength,version:[]byte{0,0,0,1},
		desConnIdLength: uint8(len(DesConnectionID)),desConnId:DesConnectionID,
		sourceConnLength:0, sourceConnId:[]byte{},tokenLength:[]byte{0},token:[]byte{},packetNumber:pn,
		payload:[]byte{1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,12,3,4,5,6,7,8},
	}
}






