package udp

import (
	"crypto/sha256"
	"github.com/chuccp/utils/log"
	"golang.org/x/crypto/hkdf"
	"testing"
	"time"
)

var quicSalt = []byte{0x38, 0x76, 0x2c, 0xf7, 0xf5, 0x59, 0x34, 0xb3, 0x4d, 0x17, 0x9a, 0xe6, 0xa4, 0xc8, 0x0c, 0xad, 0xcc, 0xbb, 0x7f, 0x0a}

func TestInitial(t *testing.T) {

	cid,err:=GenerateConnectionID(8)

	log.Info(cid,err)




	initialSecret := hkdf.Extract(sha256.New,cid, quicSalt)

	log.Info(initialSecret, len(initialSecret))

	time.Sleep(time.Second)
	
}
