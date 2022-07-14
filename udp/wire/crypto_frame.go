package wire

import (
	"github.com/chuccp/utils/udp/util"
)

type CryptoFrame struct {
	Offset uint64
	Data   []byte
}

func NewCryptoFrame(data []byte) *CryptoFrame {
	return &CryptoFrame{Data: data, Offset: 0}
}

func (cryptoFrame *CryptoFrame) Write(write *util.WriteBuffer) {
	write.WriteByte(CryptoType)
	write.WriteVariableLength(uint32(cryptoFrame.Offset))
	write.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
		wr.WriteBytes(cryptoFrame.Data)
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
	cryptoFrame.Data = data
	return nil
}
func ReadCryptoFrame(read *util.ReadBuffer) (*CryptoFrame,error) {
	var cryptoFrame CryptoFrame
	return &cryptoFrame,cryptoFrame.Read(read)
}
