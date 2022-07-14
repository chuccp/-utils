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
func (s *CipherSuites) Write(write *util.WriteBuffer)  {
	for _, suite := range s.cipherSuites {
		write.WriteUint16(suite)
	}
}
func (s *CipherSuites) Read(read *util.ReadBuffer) error{

	for{
		length, err := read.ReadUint16Length()
		if err != nil {
			return err
		}
		s.cipherSuites = append(s.cipherSuites, length)
		if read.Buffered()==0{
			break
		}
	}
	return nil
}
func  ReadCipherSuites(read *util.ReadBuffer) (*CipherSuites,error){
	var cs CipherSuites
	cs.cipherSuites = make([]uint16,0)
	return &cs,cs.Read(read)
}