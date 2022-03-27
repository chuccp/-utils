package udp

import (
	"bytes"
	"github.com/chuccp/utils/io"
	"log"
)

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

func NewClientHello() *ClientHello {
	ch:=&ClientHello{ExtMap:make(map[ExtensionsType]*Extensions)}
	ch.HandshakeType = byte(typeClientHello)
	ch.Version = []byte{0x03,0x03}
	ch.Random,_ = RandId(32)
	ch.SessionLength = 0
	ch.CiphersSuites = CipgerBytes(TLS_AES_128_GCM_SHA256,TLS_AES_256_GCM_SHA384)
	ch.CompressionMethods = []byte{0}
	ch.CompressionMethodsLength=1
	return ch
}
func (c *ClientHello) SetLength(data []byte)  {
	c.Length = uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
}
func (c *ClientHello)Bytes()[]byte{
	var buff = new(bytes.Buffer)
	buff.Write(c.Version)
	buff.Write(c.Random)
	buff.WriteByte(c.SessionLength)
	buff.Write(U16B(uint16(len(c.CiphersSuites))))
	buff.Write(c.CiphersSuites)
	buff.WriteByte(1)
	buff.WriteByte(0)
	ex:=Extansions()
	buff.Write(U16B(uint16(len(ex))))
	buff.Write(ex)
	var buff1 = new(bytes.Buffer)
	buff1.WriteByte(1)
	data:=buff.Bytes()
	bytes:=U32B(uint32(len(data)))
	buff1.Write(bytes[1:4])
	buff1.Write(data)
	return buff1.Bytes()
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