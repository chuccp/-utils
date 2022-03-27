package udp

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/io"
)

func ReadVariableLength(read *io.ReadStream) (uint32, int) {
	lenType, err := read.ReadByte()
	if err != nil {
		return 0, 0
	}
	len1 := lenType >> 6
	if len1 == 0 {
		return uint32(lenType & 63), 1
	} else {
		data, err := read.ReadBytes(int(len1))
		if err != nil {
			return 0, 0
		}
		switch len1 {
		case 1:
			return uint32(lenType&63)<<8 | uint32(data[0]), 2
		case 2:
			return uint32(lenType&63)<<16 | uint32(data[0])<<8 | uint32(data[1]), 3
		default:
			return uint32(lenType&63)<<24 | uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]), 4
		}
	}
}
func ReadBytesVariableLength(data []byte) (uint32, int) {
	return ReadVariableLength(io.NewReadStream(bytes.NewReader(data)))
}

func VariableLengthBytes(num uint32) []byte  {
	if num<=0x3f{
		return []byte{byte(num)}
	}else if num<=0x3f_ff{
		return []byte{byte(num>>8)|0x40,byte(num)}
	}else if num<=0x3f_ff_ff{
		return []byte{byte(num>>16)|0x80,byte(num>>8),byte(num)}
	}else {
		return []byte{byte(num>>24)|0xC0,byte(num>>16),byte(num>>8),byte(num)}
	}
}
func VariableLengthBytes2(num uint32) []byte  {
	if num<=0xf{
		return []byte{byte(num)}
	}else if num<=0xff{
		return []byte{byte(num>>8)|0x40,byte(num)}
	}else {
		return []byte{0x80,byte(num>>16),byte(num>>8),byte(num)}
	}
}


func U16B(value uint16)[]byte  {
	data:=make([]byte, 2)
	binary.BigEndian.PutUint16(data,value)
	return data
}
func U32B(value uint32)[]byte  {
	data:=make([]byte, 4)
	binary.BigEndian.PutUint32(data,value)
	return data
}