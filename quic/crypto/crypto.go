package crypto

import (
	"encoding/binary"
	"golang.org/x/crypto/chacha20"
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