package udp

import "testing"

func TestUn_Packet(t *testing.T) {

	data:=[]byte{128 ,255, 0, 0 ,1 ,16, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	UnPacket(data)

}
