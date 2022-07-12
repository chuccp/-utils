package tls

import "github.com/chuccp/utils/udp/util"

const (
	TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256       uint16 = 0xc02b
	TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384       uint16 = 0xc02c
)

//var CipherSuites = []uint16{TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384}

type CipherSuites struct {
	cipherSuites []uint16
}

func NewCipherSuites()*CipherSuites  {
	var data = []uint16{TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384}
	return &CipherSuites{cipherSuites:data}
}
func (s *CipherSuites) Bytes(write *util.WriteBuffer)  {
	for _, suite := range s.cipherSuites {
		write.WriteUint16(suite)
	}
}

