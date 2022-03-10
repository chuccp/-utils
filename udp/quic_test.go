package udp

import (
	"testing"
	"time"
)


func TestInitial(t *testing.T) {

	cid,_:=GenerateConnectionID(8)



	t.Log(computeSecrets(cid))
	time.Sleep(time.Second)
	
}
