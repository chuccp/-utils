package udp

import (
	"bytes"
	"encoding/binary"
)

type PacketWriteBuffer struct {
	buffer *bytes.Buffer
}

func NewPacketWriteBuffer() *PacketWriteBuffer {
	return &PacketWriteBuffer{buffer: new(bytes.Buffer)}
}
func (pb *PacketWriteBuffer) WriteByte(b byte) {
	pb.buffer.WriteByte(b)
}
func (pb *PacketWriteBuffer) WriteBytes(bs []byte) {
	if len(bs)>0{
		pb.buffer.Write(bs)
	}
}
func (pb *PacketWriteBuffer) WriteUint32(len uint8,num uint32) {
	v:=[]byte{0,0,0,0}
	binary.LittleEndian.PutUint32(v,num)
	pb.buffer.Write(v[4-len:4])
}
func (pb *PacketWriteBuffer) WriteUint64(len uint8,num uint64) {
	v:=[]byte{0,0,0,0,0,0,0,0}
	binary.LittleEndian.PutUint64(v,num)
	pb.buffer.Write(v[8-len:8])
}
func (pb *PacketWriteBuffer) WriteVariableLength(len uint32) {
	pb.buffer.Write(VariableLengthToBytes(len))
}
func (pb *PacketWriteBuffer) Bytes()[]byte {
	return  pb.buffer.Bytes()
}
