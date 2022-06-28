package tls13

import (
	"crypto/sha256"
	"golang.org/x/crypto/hkdf"
	"testing"
)

func TestName(t *testing.T) {

	key:=[]byte{1,2,3,4,5,6,7,8,9}
	value:=[]byte{1,2,3,4,5,6,7,8,9}
	hkdf.Extract(sha256.New,key,value)
}
