package udp

import (
	"encoding/binary"
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
	log2 "log"
	"testing"
	"time"
)


func TestRaw(t *testing.T) {

	fi, err := file.NewFile("C:\\Users\\cooge\\Documents\\quic\\Initial4.bin")
	if err == nil {
		file, err1 := fi.ToRawFile()
		if err1 == nil {
			read := io.NewReadStream(file)
			buff := NewBuffer()
			err1 = buff.readPack(read)
			if err1 == nil {
				log.Info("$$$$$",ConnectionID(buff.data))
				header := parseLongHeader(buff.data)
				if header.packetType==PacketTypeInitial{
					origPNBytes := make([]byte, 4)

					t.Log(header.parsedLen)

					copy(origPNBytes, buff.data[header.parsedLen:header.parsedLen+4])

					param1 :=buff.data[header.parsedLen+4:header.parsedLen+4+16]
					param2 := &buff.data[0]
					param3 := buff.data[header.parsedLen:header.parsedLen+4]
					log.Info("~~~~~~~~desConnId:",header.desConnId)
					 aead:=NewInitialAEAD(header.desConnId)
					 //b,a,nonce:
					mask:=make([]byte, aead.block.BlockSize())
					aead.block.Encrypt(mask,param1)
					log.Info("~~~~~~~~",*param2)
					*param2 ^= mask[0] & 0xf
					for i := range param3 {
						param3[i] ^= mask[i+1]
					}


					log.Info("~~~~@@@@~~~~",*param2)

					pageNum:=int(buff.data[0]&0x3)+1

					exlen:=header.parsedLen+pageNum
					nonceBuf:=make([]byte, aead.aead.NonceSize())
					binary.BigEndian.PutUint64(nonceBuf[len(nonceBuf)-8:], 0)

					log.Info("***",aead.iv)

					for i, b := range nonceBuf[len(nonceBuf)-8:] {
						aead.iv[4+i] ^= b
					}
					log.Info("***",aead.iv)
					copy(buff.data[header.parsedLen+pageNum:header.parsedLen+4],origPNBytes[pageNum:4] )

					log.Info("====",aead.iv)
					ciphertext:=buff.data[exlen:buff.len]

					log.Info("ciphertext:",len(ciphertext),"@@@@@",header.length)

					additionalData:=buff.data[:exlen]
					log.Info("====",ciphertext)
					log.Info("====",additionalData)

					data,err:=aead.aead.Open([]byte{},aead.iv,ciphertext,additionalData)

					log.Info("=============",err)

					log.Info("data:",len(data),"@@@@@",header.length,len(additionalData))

					log2.Println(ConnectionID(data),err)
					if err==nil{
						buff,_:=ParseFrame(data)
						ParseCrypto(buff)
					}
				}
			}
		}
	}

	time.Sleep(time.Second * 2)
}
