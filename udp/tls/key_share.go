package tls

import "github.com/chuccp/utils/udp/util"

type GroupType uint16
const (
	x25519 GroupType = 29
)

type ClientKeyShare struct {
	Group        GroupType
	KeyExchanges []byte
}

func NewKeyShare(groupType GroupType,KeyExchanges []byte) *ClientKeyShare {
	return &ClientKeyShare{Group: groupType, KeyExchanges: KeyExchanges}
}
func (kse *ClientKeyShare) Write(write *util.WriteBuffer) {
	write.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
		wr.WriteUint16(uint16(kse.Group))
		wr.WriteUint16(uint16(len(kse.KeyExchanges)))
		wr.WriteBytes(kse.KeyExchanges)
	})
}
func (kse *ClientKeyShare) Read(read *util.ReadBuffer) error {

	_, bytes, err := read.ReadU16LengthBytes()
	if err != nil {
		return err
	}
	rd := util.NewReadBuffer(bytes)
	u16, err := rd.ReadUint16Length()
	if err != nil {
		return err
	}
	kse.Group = GroupType(u16)
	_, kse.KeyExchanges, err = rd.ReadU16LengthBytes()
	if err != nil {
		return err
	}
	return nil
}
