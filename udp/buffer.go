package udp

import "io"

type ByteCount int64
const MaxPacketBufferSize ByteCount = 1452

type buffer struct {
	data []byte
	len uint64
}

func NewBuffer() *buffer {
	data:=make([]byte,MaxPacketBufferSize)
	return &buffer{data: data}
}
func (b *buffer) readPack(reader io.Reader) (err error) {
	len,err :=reader.Read(b.data)
	b.len = uint64(len)
	return err
}

