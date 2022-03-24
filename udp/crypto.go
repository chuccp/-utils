package udp

import (
	"bytes"
	"github.com/chuccp/utils/io"
	"log"
)

type messageType uint8

// TLS handshake message types.
const (
	typeClientHello         messageType = 1
	typeServerHello         messageType = 2
	typeNewSessionTicket    messageType = 4
	typeEncryptedExtensions messageType = 8
	typeCertificate         messageType = 11
	typeCertificateRequest  messageType = 13
	typeCertificateVerify   messageType = 15
	typeFinished            messageType = 20
)

type ClientHello struct {
	Length        uint32
	Version       []byte
	Random        []byte
	SessionLength uint8
	SessionId     []byte
	CipherSuiteLength uint16
	Ciphers []byte
	CompressionMethodsLength  uint8
	CompressionMethods []byte
	ExtensionsLength uint16
	ExtMap map[ExtensionsType]*Extensions
}

func NewClientHello() *ClientHello {
	return &ClientHello{ExtMap:make(map[ExtensionsType]*Extensions)}
}
type  ExtensionsType uint16
type Extensions struct {
	ExtensionsType ExtensionsType
	ExtensionsLength  uint16
	ExtensionsData []byte
}

func (c *ClientHello) SetLength(data []byte)  {
	c.Length = uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
}

func ParseCrypto(buff *bytes.Buffer) error{
	log.Println("@@@",ConnectionID(buff.Bytes()))
	read:=io.NewReadStream(buff)
	type_, err :=read.ReadByte()
	if err!=nil{
		return err
	}
	switch messageType(type_) {
	case typeClientHello:
		ParseClientHello(read)
	}
	return nil
}
func ParseClientHello(read *io.ReadStream)  {
	  var ch  = NewClientHello()
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

		  ch.CipherSuiteLength= uint16(cipherLength[0] )<<8|  uint16(cipherLength[1])
		  ch.Ciphers,_=read.ReadBytes(int(ch.CipherSuiteLength))
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