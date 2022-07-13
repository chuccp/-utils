package wire

import (
	"github.com/chuccp/utils/udp/util"
)

type CryptoFrame struct {
	Offset uint64
	data   []byte
}

func NewCryptoFrame(data []byte) *CryptoFrame {
	return &CryptoFrame{data: data, Offset: 0}
}

func (cryptoFrame *CryptoFrame) Write(write *util.WriteBuffer) {
	write.WriteByte(CryptoType)
	write.WriteVariableLength(uint32(cryptoFrame.Offset))
	write.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
		wr.WriteBytes(cryptoFrame.data)
	})
}
func (cryptoFrame *CryptoFrame) Read(read *util.ReadBuffer) error {
	length, err := read.ReadVariableLength()
	if err != nil {
		return err
	}
	cryptoFrame.Offset = uint64(length)
	_, data, err := read.ReadVariableLengthBytes()
	if err != nil {
		return err
	}
	cryptoFrame.data = data
	return nil
}
