package wire

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/udp/util"
	"log"
	"sync"
)

type   HandShakeProgress uint8
const (
	WaitRetry  HandShakeProgress = iota

)



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
func initAEAD(key [16]byte) cipher.AEAD {
	aes, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}
	return aead
}
var (
	retryAEAD = initAEAD([16]byte{0xbe, 0x0c, 0x69, 0x0b, 0x9f, 0x66, 0x57, 0x5a, 0x1d, 0x76, 0x6b, 0x54, 0xe3, 0x68, 0xc8, 0x4e})
	retryNonce    = [12]byte{0x46, 0x15, 0x99, 0xd3, 0x5d, 0x63, 0x2b, 0xf2, 0x23, 0x98, 0x25, 0xbb}
)

var (
	retryMutex    sync.Mutex
)

func GetRetryIntegrityTag(retry []byte, origDestConnID ConnectionID) *[16]byte {
	retryMutex.Lock()
	var  retryBuf = new(bytes.Buffer)
	retryBuf.WriteByte(uint8(origDestConnID.Len()))
	retryBuf.Write(origDestConnID.Bytes())
	retryBuf.Write(retry)
	var tag [16]byte
	var sealed = retryAEAD.Seal(tag[:0], retryNonce[:], nil, retryBuf.Bytes())
	if len(sealed) != 16 {
		panic(fmt.Sprintf("unexpected Retry integrity tag length: %d", len(sealed)))
	}
	retryMutex.Unlock()
	return &tag
}

func HandleRetryPacket(data []byte,desConnId ConnectionID) bool  {
	l:= len(data)
	retryToken := data[:l-16]

	tagV:=GetRetryIntegrityTag(retryToken,desConnId.Bytes())
	tag :=data[l-16:]
	if bytes.Equal(tagV[:],tag){
		return true
	}
	return false
}
type  ExtensionsType uint16
type Extensions struct {
	ExtensionsType ExtensionsType
	ExtensionsLength  uint16
	ExtensionsData []byte
}

type ClientHello struct {
	HandshakeType            byte
	Length                   uint32
	Version                  []byte
	Random                   []byte
	SessionLength            uint8
	SessionId                []byte
	CiphersSuiteLength       uint16
	CiphersSuites            []byte
	CompressionMethodsLength uint8
	CompressionMethods       []byte
	ExtensionsLength         uint16
	ExtMap                   map[ExtensionsType]*Extensions
}

func NewClientHello(random []byte) *ClientHello {
	ch:=&ClientHello{ExtMap:make(map[ExtensionsType]*Extensions)}
	ch.HandshakeType = byte(typeClientHello)
	ch.Version = []byte{0x03,0x03}
	ch.Random = random[:]
	ch.SessionLength = 0
	ch.CiphersSuites = CipgerBytes(
		TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		TLS_RSA_WITH_AES_128_GCM_SHA256,
		TLS_RSA_WITH_AES_256_GCM_SHA384,
		TLS_RSA_WITH_AES_128_CBC_SHA,
		TLS_RSA_WITH_AES_256_CBC_SHA,
		TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		TLS_AES_128_GCM_SHA256,
		TLS_AES_256_GCM_SHA384,
		TLS_CHACHA20_POLY1305_SHA256)
	ch.CompressionMethods = []byte{0}
	ch.CompressionMethodsLength=1
	return ch
}
func (c *ClientHello) SetLength(data []byte)  {
	c.Length = uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
}
func (c *ClientHello)Bytes(nextProtos []string)[]byte{
	var buff = new(bytes.Buffer)
	buff.Write(c.Version)
	buff.Write(c.Random)
	buff.WriteByte(c.SessionLength)
	buff.Write(util.U16B(uint16(len(c.CiphersSuites))))
	buff.Write(c.CiphersSuites)
	buff.WriteByte(1)
	buff.WriteByte(0)
	ex:=Extansions(nextProtos)
	buff.Write(util.U16B(uint16(len(ex))))
	buff.Write(ex)
	var buff1 = new(bytes.Buffer)
	buff1.WriteByte(1)
	data:=buff.Bytes()
	bytes:=util.U32B(uint32(len(data)))
	buff1.Write(bytes[1:4])
	buff1.Write(data)
	return buff1.Bytes()
}

func ParseClientHello(read *io.ReadStream)  {
	var ch  = new(ClientHello)
	data,err:=read.ReadBytes(3)
	if err==nil{
		ch.Length= uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
		log.Printf("%d", ch.Length)
		ch.Version,_= read.ReadBytes(2)
		ch.Random,_= read.ReadBytes(32)
		ch.SessionLength,_=read.ReadByte()
		if ch.SessionLength>0 {
			ch.SessionId,_= read.ReadBytes(int(ch.SessionLength))
		}
		cipherLength,_:=read.ReadBytes(2)

		ch.CiphersSuiteLength = uint16(cipherLength[0] )<<8|  uint16(cipherLength[1])
		ch.CiphersSuites,_=read.ReadBytes(int(ch.CiphersSuiteLength))
		log.Println("SessionLength:",ch.SessionLength)
		ch.CompressionMethodsLength,_=read.ReadByte()

		log.Println("CompressionMethodsLength:",ch.CompressionMethodsLength)
		if ch.CompressionMethodsLength>0{
			ch.CompressionMethods,_=read.ReadBytes(int(ch.CompressionMethodsLength))
			log.Println("CompressionMethods:",ch.CompressionMethods)
		}
		extensionsLength,_:=read.ReadBytes(2)
		log.Println(extensionsLength)
		ch.ExtensionsLength = uint16(extensionsLength[0] )<<8|  uint16(extensionsLength[1])

		log.Printf("%d",ch.ExtensionsLength)
		extensions,_:=read.ReadBytes(int(ch.ExtensionsLength))
		extension:=io.NewReadBytesStream(extensions)
		ParseExtension(extension, func(ext *Extensions) {
			ch.ExtMap[ext.ExtensionsType] = ext
		})


	}
}
func ParseExtension(read *io.ReadStream,f func(ext *Extensions)) {

	for{
		var ext  = &Extensions{}
		extensionType,err:=read.Read2Uint16()
		if err!=nil{
			break
		}
		ext.ExtensionsType = ExtensionsType(extensionType)
		ext.ExtensionsLength, err = read.Read2Uint16()
		log.Println(ext.ExtensionsType,"  ",ext.ExtensionsLength)
		if err!=nil{
			break
		}
		ext.ExtensionsData,err = read.ReadBytes(int(ext.ExtensionsLength))
		if err!=nil{
			break
		}
		f(ext)
	}
}
