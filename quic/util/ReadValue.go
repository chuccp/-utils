package util

import "github.com/chuccp/utils/io"

type ReadValue struct {
	stream *io.ReadStream
}

func NewReadValue(stream *io.ReadStream) *ReadValue {
	return &ReadValue{stream: stream}
}

func (read *ReadValue) ReadUint8() ([]byte, uint8, error) {
	i8, err := read.stream.ReadUint8()
	if err != nil {
		return nil, i8, err
	}
	data, err := read.stream.ReadBytes(int(i8))
	if err != nil {
		return nil, i8, err
	} else {
		return data, i8, err
	}
}

func (read *ReadValue) ReadVariableValueLength() ( uint32, error) {
	i8, err := read.readVariableLength(read.stream)
	if err != nil {
		return i8, err
	}else{
		return i8, nil
	}
	//data, err := read.stream.ReadBytes(int(i8))
	//if err != nil {
	//	return  i8, err
	//} else {
	//	return data, i8, err
	//}
}

func (read *ReadValue)readVariableLength(read1 *io.ReadStream) (uint32, error) {
	lenType, err := read1.ReadByte()
	if err != nil {
		return 0, err
	}
	len1 := lenType >> 6
	if len1 == 0 {
		return uint32(lenType & 63), nil
	} else {
		data, err := read1.ReadBytes(int(len1))
		if err != nil {
			return 0, err
		}
		switch len1 {
		case 1:
			return uint32(lenType&63)<<8 | uint32(data[0]), nil
		case 2:
			return uint32(lenType&63)<<16 | uint32(data[0])<<8 | uint32(data[1]), nil
		default:
			return uint32(lenType&63)<<24 | uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]), nil
		}
	}
}
