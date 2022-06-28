package crypto

import (
	"math/big"
	"testing"
)

func TestQQQQ(t *testing.T) {

	ecdh0, err := NewECDH()
	if err != nil {
		t.Log(ecdh0)
	}
	ecdh1, err := NewECDH()
	if err != nil {
		t.Log(ecdh1)
	}

	t.Log(ecdh0.GetAESKey(ecdh1.public))
	t.Log(ecdh1.GetAESKey(ecdh0.public))
}
func TestName(t *testing.T) {

	t.Log(big.NewInt(1))

}
