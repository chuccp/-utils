package udp

import (
	"log"
	"os"
	"testing"
)

func TestUn_Packet(t *testing.T) {

	data, err := os.ReadFile("data.bb")
	if err != nil {
		return
	}
	log.Print(data)
	var longHeader LongHeader
	err = UnPacket(data,&longHeader)
	if err != nil {
		return 
	}

	if longHeader.IsLongHeader{
		if longHeader.LongPacketType==packetTypeInitial{
			UnPacketInitialPayload(&longHeader)
		}
	}
}
