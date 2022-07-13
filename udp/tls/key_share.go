package tls

import "github.com/chuccp/utils/udp/util"




type KeyShare struct {
	Group []byte
	KeyExchanges []byte
}

func NewKeyShare(KeyExchanges []byte) *KeyShare {
	group := []byte{0x00,0x0d}
	return &KeyShare{Group:group,KeyExchanges:KeyExchanges}
}
func (kse *KeyShare) Write(write *util.WriteBuffer)  {
	write.WriteBytes(kse.Group)
	write.WriteUint16(uint16(len(kse.KeyExchanges)))
	write.WriteBytes(kse.KeyExchanges)
}