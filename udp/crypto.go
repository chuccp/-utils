package udp

import (
	"bytes"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
)

type messageType uint8

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

type Crypto struct {
	Offset uint32
	Length uint32
	data []byte
}

func  NewCrypto(Offset uint32,data []byte)*Crypto  {
	log.Info("NewCrypto:", len(data),ConnectionID(data))
	return &Crypto{Offset:Offset,Length: uint32(len(data)),data:data}
}

func (c *Crypto) Bytes()[]byte  {
	var buff = new(bytes.Buffer)
	buff.WriteByte(0x06)
	buff.Write(VariableLengthBytes(c.Offset))
	buff.Write(VariableLengthBytes(c.Length))
	buff.Write(c.data)
	return buff.Bytes()
}


func ParseCrypto(buff *bytes.Buffer) error{
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
