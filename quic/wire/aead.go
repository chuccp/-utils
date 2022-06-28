package wire

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/hkdf"
)

var (
	//quicSaltOld = []byte{0xaf, 0xbf, 0xec, 0x28, 0x99, 0x93, 0xd2, 0x4c, 0x9e, 0x97, 0x86, 0xf1, 0x9c, 0x61, 0x11, 0xe0, 0x43, 0x90, 0xa8, 0x99}
	quicSalt = []byte{0x38, 0x76, 0x2c, 0xf7, 0xf5, 0x59, 0x34, 0xb3, 0x4d, 0x17, 0x9a, 0xe6, 0xa4, 0xc8, 0x0c, 0xad, 0xcc, 0xbb, 0x7f, 0x0a}
)

func hkdfExpandLabel(hash crypto.Hash, secret, context []byte, label string, length int) []byte {
	b := make([]byte, 3, 3+6+len(label)+1+len(context))
	binary.BigEndian.PutUint16(b, uint16(length))
	b[2] = uint8(6 + len(label))
	b = append(b, []byte("tls13 ")...)
	b = append(b, []byte(label)...)
	b = b[:3+6+len(label)+1]
	b[3+6+len(label)] = uint8(len(context))
	b = append(b, context...)
	out := make([]byte, length)
	n, err := hkdf.Expand(hash.New, secret, b).Read(out)
	if err != nil || n != length {
		panic("quic: HKDF-Expand-Label invocation failed unexpectedly")
	}
	return out
}

func newAESHeaderProtector(trafficSecret []byte) cipher.Block {
	hpKey := hkdfExpandLabel(crypto.SHA256, trafficSecret, []byte{}, "quic hp", 16)
	block, err := aes.NewCipher(hpKey)
	if err != nil {
		panic(fmt.Sprintf("error creating new AES cipher: %s", err))
	}
	return block
}

func aeadAESGCM(key []byte) cipher.AEAD {
	aes, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}
	return aead
}

func computeSecrets(connID []byte) (clientSecret, serverSecret []byte) {
	initialSecret := hkdf.Extract(sha256.New, connID, quicSalt)
	clientSecret = hkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "client in", crypto.SHA256.Size())
	serverSecret = hkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "server in", crypto.SHA256.Size())
	return
}

func computeInitialKeyAndIV(secret []byte) (key, iv []byte) {
	key = hkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic key", 16)
	iv = hkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic iv", 12)
	return
}

type AEAD struct {
	block cipher.Block
	aead  cipher.AEAD
	key   []byte
	iv    []byte
}



func (a *AEAD)NonceSize() int{
	return 0
}
func (a *AEAD)Overhead() int{
	return 0
}


func (a *AEAD)Seal(dst, nonce, plaintext, additionalData []byte) []byte{
	return nil
}


func (a *AEAD)Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error){
	return a.aead.Open(dst,nonce,ciphertext,additionalData)
}

func (a *AEAD)Encrypt(dst, src []byte){

	a.block.Encrypt(dst,src)
}


func (a *AEAD)Decrypt(dst, src []byte){
	a.block.Decrypt(dst,src)
}

func NewInitialAEAD(connID []byte, isClient bool) *AEAD {
	clientSecret, serverSecret := computeSecrets(connID)
	var mySecret, _ []byte
	if isClient {
		mySecret = clientSecret
		_ = serverSecret
	} else {
		_ = clientSecret
		mySecret = serverSecret
	}
	block := newAESHeaderProtector(mySecret)
	myKey, myIV := computeInitialKeyAndIV(mySecret)
	aead := aeadAESGCM(myKey)
	return &AEAD{block: block, aead: aead, key: myKey, iv: myIV}
}
