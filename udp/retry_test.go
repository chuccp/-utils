package udp

import (
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/udp/wire"
	"testing"
)

func TestRetry(t *testing.T) {
	fi, err := file.NewFile("C:\\Users\\cooge\\Documents\\quic\\retry.bin")
	if err==nil{
		data:=make([]byte,MaxPacketBufferSize)
		n,err:=fi.ReadBytes(data)
		if err==nil{
			wire.ParsePacket(data[:n])
		}
	}
}
func TestRetryyy(t *testing.T)  {
	t.Log(wire.GetRetryIntegrityTag([]byte("foobar"),wire.ConnectionID{1, 2, 3, 4}))
}