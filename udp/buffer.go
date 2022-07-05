package udp

import (
	"bytes"
	"encoding/binary"
)

type PacketBuffer struct {
	buffer *bytes.Buffer
}

func NewPacketBuffer() *PacketBuffer {
	return &PacketBuffer{buffer: new(bytes.Buffer)}
}
func (pb *PacketBuffer) WriteByte(b byte) {
	pb.buffer.WriteByte(b)
}
func (pb *PacketBuffer) WriteBytes(bs []byte) {
	if len(bs)>0{
		pb.buffer.Write(bs)
	}
}
func (pb *PacketBuffer) WriteUint32(len uint8,num uint32) {
	v:=[]byte{0,0,0,0}
	binary.LittleEndian.PutUint32(v,num)
	pb.buffer.Write(v[4-len:4])
}
func (pb *PacketBuffer) WriteUint64(len uint8,num uint64) {
	v:=[]byte{0,0,0,0,0,0,0,0}
	binary.LittleEndian.PutUint64(v,num)
	pb.buffer.Write(v[8-len:8])
}
func (pb *PacketBuffer) WriteVariableLength(len uint32) {
	pb.buffer.Write(VariableLengthToBytes(len))
}
func (pb *PacketBuffer) Bytes()[]byte {
	return  pb.buffer.Bytes()
}
