package util

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/io"
	"log"
)

type BufferWrite interface {
	Write(write *WriteBuffer)
}
type BufferRead interface {
	Read(read *ReadBuffer)error
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
	log.Print(uint32(len(data)))
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
	binary.BigEndian.PutUint32(v, num)
	pb.buffer.Write(v[4-len : 4])
}
func (pb *WriteBuffer) WriteUint16(num uint16) {
	v := []byte{0, 0}
	binary.BigEndian.PutUint16(v, num)
	pb.buffer.Write(v)
}
func (pb *WriteBuffer) WriteLenUint64(len uint8, num uint64) {
	v := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(v, num)
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
	binary.BigEndian.PutUint32(v, u)
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
func (prb *ReadBuffer) ReadByteBuff(f func(byte, *ReadBuffer) error) error {
	readByte, err := prb.ReadByte()
	if err != nil {
		return err
	}
	return f(readByte, prb)
}
func (prb *ReadBuffer) ReadBytesBuff(num uint8, f func([]byte, *ReadBuffer) error) error {
	readBytes, err := prb.ReadBytes(int(num))
	if err != nil {
		return err
	}
	return f(readBytes, prb)
}
func (prb *ReadBuffer) ReadUInt32Buff(f func(uint32, *ReadBuffer) error) error {
	u32, err := prb.Read4U32()
	if err != nil {
		return err
	}
	return f(u32, prb)
}
func (prb *ReadBuffer) ReadUint8LengthBytesBuff(f func([]byte, *ReadBuffer) error) error {
	u8, err := prb.ReadByte()
	if err != nil {
		return err
	}
	if u8 == 0 {
		return f([]byte{}, prb)
	}
	readBytes, err := prb.ReadBytes(int(u8))
	if err != nil {
		return err
	}
	return f(readBytes, prb)
}
func (prb *ReadBuffer) ReadUint8LengthU32Buff(f func(uint32, *ReadBuffer) error) error {
	u8, err := prb.ReadByte()
	if err != nil {
		return err
	}
	if u8 == 0 {
		return f(0, prb)
	}
	readBytes, err := prb.ReadU8LengthU32(u8)
	if err != nil {
		return err
	}
	return f(readBytes, prb)
}

func (prb *ReadBuffer) ReadVariableLengthBuff(f func(uint32, *ReadBuffer) error) error {
	u32, err := prb.ReadVariableLength()
	if err!=nil{
		return err
	}
	return f(u32, prb)
}


func (prb *ReadBuffer) ReadVariableLengthBytesBuff(f func([]byte, *ReadBuffer) error) error {
	num, err := prb.ReadVariableLength()
	if err != nil {
		return err
	}
	if num == 0 {

		return f([]byte{}, prb)
	}
	readBytes, err := prb.ReadBytes(int(num))
	if err != nil {
		return err
	}
	return f(readBytes, prb)
}

func (prb *ReadBuffer) ReadU32Bytes(u32 uint32) ([]byte, error) {
	return prb.readStream.ReadUintBytes(u32)
}

func (prb *ReadBuffer) ReadU8LengthU32(u8 uint8) (uint32, error) {
	if u8==0 {
		return 0,nil
	}

	data, err := prb.readStream.ReadBytes(int(u8))
	if err != nil {
		return 0, err
	}
	switch u8 {
	case 1:
		return  uint32(data[0]),nil
	case 2:
		{
			return	uint32(data[0])<<8|uint32(data[1]),nil
		}
	case 3:
		{
			return	uint32(data[0])<<16|uint32(data[1])<<8|uint32(data[2]),nil
		}
	case 4:
		{
			return	uint32(data[0])<<24|uint32(data[1])<<16|uint32(data[2])<<8|uint32(data[3]),nil
		}
	}


	return 0, nil
}

func (prb *ReadBuffer) ReadBytes(len int) ([]byte, error) {
	return prb.readStream.ReadBytes(len)
}
func (prb *ReadBuffer) Read4U32() (uint32, error) {
	return prb.readStream.Read4Uint32()
}
func (prb *ReadBuffer) ReadU8LengthBytes() (uint8, []byte, error) {
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
func (prb *ReadBuffer) ReadU16LengthBytes() (uint16, []byte, error) {
	u16, err := prb.readStream.Read2Uint16()
	if err != nil {
		return 0, nil, err
	}
	if u16 == 0 {
		return u16, []byte{}, nil
	}
	readBytes, err := prb.readStream.ReadBytes(int(u16))
	if err != nil {
		return u16, nil, err
	}
	return u16, readBytes, nil
}
func (prb *ReadBuffer) Offset() uint16 {
	return prb.readStream.Offset()
}
func (prb *ReadBuffer) Size() int {
	return prb.readStream.Size()
}
func (prb *ReadBuffer) Buffered() int {
	return prb.readStream.Buffered()
}

func (prb *ReadBuffer) ReadUint24Length()(uint32, error) {
	return prb.readStream.Read3Uint32()
}
func (prb *ReadBuffer) ReadUint16Length()(uint16, error) {
	return prb.readStream.Read2Uint16()
}
func (prb *ReadBuffer) ReadVariableLength() (uint32, error) {

	b, err := prb.readStream.ReadByte()
	if err != nil {
		return 0, err
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
				return 0, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<8) | uint32(b2)
		}
	case 2:
		{
			readBytes, err := prb.readStream.ReadBytes(2)
			if err != nil {
				return 0, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<16) | (uint32(readBytes[0]) << 8) | uint32(readBytes[1])
		}
	case 3:
		{
			readBytes, err := prb.readStream.ReadBytes(3)
			if err != nil {
				return 0, err
			}
			b = b & 0x3F
			num = (num | uint32(b)<<24) | (uint32(readBytes[0]) << 16) | uint32(readBytes[1])<<8 | uint32(readBytes[2])
		}
	}
	return num, nil
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


