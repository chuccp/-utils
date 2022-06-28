package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/curve25519"
	io2 "io"
)

func Chacha20Mask(sample, mask []byte) {
	var key   [32]byte
	c, err := chacha20.NewUnauthenticatedCipher(key[:], sample[4:])
	if err != nil {
		panic(err)
	}
	c.SetCounter(binary.LittleEndian.Uint32(sample[:4]))
	c.XORKeyStream(mask, mask)
}

type ECDH struct {
	private []byte
	public []byte
}

func NewECDH() (*ECDH,error) {
	var echd ECDH
	echd.private = make([]byte, curve25519.ScalarSize)
	_, err := io2.ReadFull(rand.Reader,echd.private)
	if err!=nil{
		return nil, err
	}
	echd.public, err = curve25519.X25519(echd.private,curve25519.Basepoint)
	if err!=nil{
		return nil, err
	}
	return &echd, nil
}
func (e *ECDH) GetPrivate() []byte {
	return e.private
}
func (e *ECDH) GetPublic() []byte {
	return e.public
}
func (e *ECDH) GetAESKey(public []byte)([]byte,error)  {
	return curve25519.X25519(e.private,public)
}





