package udp

import (
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/tls"
	"github.com/chuccp/utils/udp/util"
	"github.com/chuccp/utils/udp/wire"
	"testing"
)

func TestInitial(t *testing.T) {

	key := []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	//rand.Read(key)
	sc := config.NewSendConfig(key)

	wb := util.NewWriteBuffer()

	clientHello := tls.NewClientHello(sc)
	clientHelloWb := util.NewWriteBuffer()
	clientHello.Write(clientHelloWb)
	cryptoFrame := wire.NewCryptoFrame(clientHelloWb.Bytes())
	cryptoFrameWb := util.NewWriteBuffer()
	cryptoFrame.Write(cryptoFrameWb)
	head := NewLongHeader(packetTypeInitial, cryptoFrameWb.Bytes(), sc)
	head.Write(wb)
	newFile, err := file.NewFile("data.bb")
	if err != nil {
		t.Log(err)
		return
	}
	err = newFile.WriteBytes(wb.Bytes())
	if err != nil {
		t.Log(err)
		return
	}

}
func TestInitialVersionNumber(t *testing.T) {

	t.Log(util.Version1.ToBytes())
}
