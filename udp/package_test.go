package udp

import "testing"

func TestLongHeaderPacket(t *testing.T) {
	connId:=[]byte{0,1,2,3,4,5,6,7,0,1,2,3,4,5,6,7}
	pack:=Initial(connId)
	t.Log(pack.Bytes())
}
