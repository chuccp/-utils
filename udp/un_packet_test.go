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
	UnPacket(data)



}
