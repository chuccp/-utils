package udp

import (
	"testing"
)

func TestInitial(t *testing.T) {

	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1}
	//rand.Read(key)
	sc := NewSendConfig(key)


	head:=NewLongHeader(packetTypeInitial,[]byte{0,0,0,0,0,0,0,0},sc)
	var packetBuffer = NewPacketWriteBuffer()
	Packet(head,packetBuffer)
	t.Log(packetBuffer.Bytes())

}
func TestInitialVersionNumber(t *testing.T)  {

	t.Log(Version1.ToBytes())
}