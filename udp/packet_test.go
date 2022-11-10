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
	sc := config.CreateSendConfig()
	clientHello := tls.CreateClientHello(sc)
	cryptoFrame := wire.CreateCryptoFrame(clientHello.Bytes())
	head := wire.CreateInitialLongHeader(cryptoFrame.Bytes(), sc)
	newFile, err := file.NewFile("data001.bb")
	if err != nil {
		t.Log(err)
		return
	}
	err = newFile.WriteBytes(head.Bytes())
	if err != nil {
		t.Log(err)
		return
	}

}
func TestInitialVersionNumber(t *testing.T) {

	t.Log(util.Version1.ToBytes())
}
