package udp

import (
	"fmt"
	"github.com/chuccp/utils/file"
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/tls"
	"github.com/chuccp/utils/udp/util"
	"github.com/chuccp/utils/udp/wire"
	"math/rand"
	"testing"
)

func TestInitial(t *testing.T) {

	key := []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	rand.Read(key)
	sc := config.NewSendConfig(key)

	clientHello := tls.CreateClientHello(sc)
	fmt.Println(clientHello.Bytes())
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
