package tls

import "github.com/chuccp/utils/udp/util"




type ClientKeyShare struct {
	Group []byte
	KeyExchanges []byte
}

func NewKeyShare(KeyExchanges []byte) *ClientKeyShare {
	group := []byte{0x00,0x1d}
	return &ClientKeyShare{Group: group,KeyExchanges:KeyExchanges}
}
func (kse *ClientKeyShare) Write(write *util.WriteBuffer)  {
	write.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
		wr.WriteBytes(kse.Group)
		wr.WriteUint16(uint16(len(kse.KeyExchanges)))
		wr.WriteBytes(kse.KeyExchanges)
	})
}