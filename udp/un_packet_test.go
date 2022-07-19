package udp

import (
	"os"
	"testing"
)

func TestUn_Packet(t *testing.T) {

	data, err := os.ReadFile("data.bb")
	if err != nil {
		return
	}
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
