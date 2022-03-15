package udp

import (
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
	str "github.com/chuccp/utils/string"
	"testing"
	"time"
)


func TestRaw(t *testing.T) {

	fi, err := file.NewFile("C:\\Users\\cooge\\Documents\\quic\\Initial.bin")
	if err == nil {
		file, err1 := fi.ToRawFile()
		if err1 == nil {
			read := io.NewReadStream(file)
			buff := NewBuffer()
			err1 = buff.readPack(read)
			if err1 == nil {
				header := parseLongHeader(buff.data)
				if header.packetType==PacketTypeInitial{
					origPNBytes := make([]byte, 4)
					copy(origPNBytes, buff.data[header.parsedLen:header.parsedLen+4])

					param1 :=buff.data[header.parsedLen+4:header.parsedLen+4+16]
					param2 := &buff.data[0]
					param3 := buff.data[header.parsedLen:header.parsedLen+4]
					log.Info("~~~~~~~~",str.BytesToHex(header.desConnId))
					b,mask:=NewInitialAEAD(header.desConnId)
					b.Encrypt(mask,param1)
					log.Info("~~~~~~~~",*param2)
					*param2 ^= mask[0] & 0xf
					for i := range param3 {
						param3[i] ^= mask[i+1]
					}

					log.Info("~~~~~~~~",*param2)


				}
			}
		}
	}

	time.Sleep(time.Second * 2)
}
