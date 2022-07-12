package udp

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/tls"
	"github.com/chuccp/utils/udp/util"
	"github.com/chuccp/utils/udp/wire"
	"log"
	"testing"
)

func TestInitial(t *testing.T) {

	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1}
	//rand.Read(key)
	sc := config.NewSendConfig(key)

	wb:= util.NewWriteBuffer()

	cryptoFrame :=wire.NewCryptoFrame(tls.NewClientHello(sc),0)
	head:=NewLongHeader(packetTypeInitial,cryptoFrame,sc)
	head.Bytes(wb)

	log.Print(wb.Bytes())


}
func TestInitialVersionNumber(t *testing.T)  {

	t.Log(util.Version1.ToBytes())
}