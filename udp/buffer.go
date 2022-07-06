package udp

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/io"
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
	if len(bs) > 0 {
		pb.buffer.Write(bs)
	}
}
func (pb *PacketWriteBuffer) WriteUint32(len uint8, num uint32) {
	v := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(v, num)
	pb.buffer.Write(v[4-len : 4])
}
func (pb *PacketWriteBuffer) WriteUint64(len uint8, num uint64) {
	v := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(v, num)
	pb.buffer.Write(v[8-len : 8])
}
func (pb *PacketWriteBuffer) WriteVariableLength(len uint32) {
	pb.buffer.Write(VariableLengthToBytes(len))
}
func (pb *PacketWriteBuffer) Bytes() []byte {
	return pb.buffer.Bytes()
}

type PacketReadBuffer struct {
	readStream *io.ReadStream
}

func NewPacketReadBuffer(data []byte) *PacketReadBuffer {
	return &PacketReadBuffer{readStream: io.NewReadBytesStream(data)}
}
func (prb *PacketReadBuffer) ReadByte() (byte, error) {
	return prb.readStream.ReadByte()
}
func (prb *PacketReadBuffer) ReadU32Bytes(u32 uint32) ([]byte, error) {
	return prb.readStream.ReadUintBytes(u32)
}
func (prb *PacketReadBuffer) ReadBytes(len int) ([]byte, error) {
	return prb.readStream.ReadBytes(len)
}
func (prb *PacketReadBuffer) Read4U32() (uint32, error) {
	return prb.readStream.Read4Uint32()
}
func (prb *PacketReadBuffer) ReadU8Bytes() (uint8, []byte, error) {
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
func (prb *PacketReadBuffer) Offset()uint16{
	return prb.readStream.Offset()
}

func (prb *PacketReadBuffer) ReadVariableLengthBytes() (uint32, []byte, error) {

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
