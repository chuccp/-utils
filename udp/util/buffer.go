package util

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/io"
)

type Buffer interface {
	Bytes(write *WriteBuffer)
}

type WriteBuffer struct {
	buffer *bytes.Buffer
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{buffer: new(bytes.Buffer)}
}
func (pb *WriteBuffer) WriteByte(b byte) {
	pb.buffer.WriteByte(b)
}
func (pb *WriteBuffer) WriteUint8LengthBytes(bs []byte) {
	ln := len(bs)
	if ln > 0 {
		pb.WriteByte(byte(ln))
		pb.WriteBytes(bs)
	} else {
		pb.WriteByte(0)
	}

}
func (pb *WriteBuffer) WriteVariableLengthBytes(data []byte) {
	ln := uint32(len(data))
	pb.WriteVariableLength(ln)
	if ln > 0 {
		pb.WriteBytes(data)
	}
}

func (pb *WriteBuffer) WriteVariableLengthBuff(f func(write *WriteBuffer)) {
	wb := NewWriteBuffer()
	f(wb)
	data := wb.Bytes()
	pb.WriteVariableLength(uint32(len(data)))
	pb.WriteBytes(data)
}

func (pb *WriteBuffer) WriteUint16LengthBuff(f func(write *WriteBuffer)) {
	wb := NewWriteBuffer()
	f(wb)
	data := wb.Bytes()
	pb.WriteUint16(uint16(len(data)))
	pb.WriteBytes(data)
}
func (pb *WriteBuffer) WriteUint24LengthBuff(f func(write *WriteBuffer)) {
	wb := NewWriteBuffer()
	f(wb)
	data := wb.Bytes()
	pb.WriteUint24(uint32(len(data)))
	pb.WriteBytes(data)
}

func (pb *WriteBuffer) WriteBytes(bs []byte) {
	if len(bs) > 0 {
		pb.buffer.Write(bs)
	}
}
func (pb *WriteBuffer) WriteUint32(len uint8, num uint32) {
	v := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(v, num)
	pb.buffer.Write(v[4-len : 4])
}
func (pb *WriteBuffer) WriteUint16(num uint16) {
	v := []byte{0, 0}
	binary.LittleEndian.PutUint16(v, num)
	pb.buffer.Write(v)
}
func (pb *WriteBuffer) WriteUint64(len uint8, num uint64) {
	v := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(v, num)
	pb.buffer.Write(v[8-len : 8])
}
func (pb *WriteBuffer) WriteVariableLength(len uint32) {
	pb.buffer.Write(VariableLengthToBytes(len))
}
func (pb *WriteBuffer) Bytes() []byte {
	return pb.buffer.Bytes()
}

func (pb *WriteBuffer) WriteUint24(u uint32) {
	v := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(v, u)
	pb.buffer.Write(v[1:4])
}

type ReadBuffer struct {
	readStream *io.ReadStream
}

func NewReadBuffer(data []byte) *ReadBuffer {
	return &ReadBuffer{readStream: io.NewReadBytesStream(data)}
}
func (prb *ReadBuffer) ReadByte() (byte, error) {
	return prb.readStream.ReadByte()
}
func (prb *ReadBuffer) ReadU32Bytes(u32 uint32) ([]byte, error) {
	return prb.readStream.ReadUintBytes(u32)
}
func (prb *ReadBuffer) ReadBytes(len int) ([]byte, error) {
	return prb.readStream.ReadBytes(len)
}
func (prb *ReadBuffer) Read4U32() (uint32, error) {
	return prb.readStream.Read4Uint32()
}
func (prb *ReadBuffer) ReadU8Bytes() (uint8, []byte, error) {
	u8, err := prb.readStream.ReadUint8()
	if err != nil {
		return 0, nil, err
	}
	if u8 == 0 {
		return u8, []byte{}, nil
	}
	readBytes, err := prb.readStream.ReadBytes(int(u8))
	if err != nil {
		return u8, nil, err
	}
	return u8, readBytes, nil
}
func (prb *ReadBuffer) Offset() uint16 {
	return prb.readStream.Offset()
}

func (prb *ReadBuffer) ReadVariableLengthBytes() (uint32, []byte, error) {

	b, err := prb.readStream.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	num := uint32(0)
	v := b >> 6
	switch v {
	case 0:
		num = uint32(b)
	case 1:
		{
			b2, err := prb.readStream.ReadByte()
			if err != nil {
				return 0, nil, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<8) | uint32(b2)
		}
	case 2:
		{
			readBytes, err := prb.readStream.ReadBytes(2)
			if err != nil {
				return 0, nil, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<16) | (uint32(readBytes[0]) << 8) | uint32(readBytes[1])
		}
	case 3:
		{
			readBytes, err := prb.readStream.ReadBytes(3)
			if err != nil {
				return 0, nil, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<24) | (uint32(readBytes[0]) << 16) | uint32(readBytes[1])<<8 | uint32(readBytes[2])
		}

	}
	readBytes, err := prb.readStream.ReadUintBytes(num)
	if err != nil {
		return 0, nil, err
	}
	return num, readBytes, nil
}
