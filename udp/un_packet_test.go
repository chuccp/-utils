package udp

import (
	"github.com/chuccp/utils/udp/wire"
	"os"
	"testing"
)

func TestUn_Packet(t *testing.T) {

	data, err := os.ReadFile("data.bb")
	if err != nil {
		return
	}
	var longHeader LongHeader
	err = UnLongHeaderPacket(data,&longHeader)
	if err != nil {
		return 
	}

	if longHeader.IsLongHeader{
		if longHeader.LongPacketType==wire.PacketTypeInitial{
			UnPacketInitialPayload(&longHeader)
		}
	}
}
