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
	UnPacket(data)



}
