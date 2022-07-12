package wire

import (
	"github.com/chuccp/utils/udp/util"
)

type CryptoFrame struct {
	Offset uint64
	Buffer util.Buffer
}

func NewCryptoFrame(Buffer util.Buffer,Offset uint64) *CryptoFrame {
	return &CryptoFrame{Buffer:Buffer,Offset:Offset}
}

func (cryptoFrame *CryptoFrame) Bytes(write *util.WriteBuffer)  {
	write.WriteByte(CryptoFrameType)
	write.WriteVariableLength(uint32(cryptoFrame.Offset))
	write.WriteVariableLengthBuff(func(wr *util.WriteBuffer) {
		cryptoFrame.Buffer.Bytes(wr)
	})
}