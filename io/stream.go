package io

import (
	"bufio"
	"bytes"
	"io"
)

type ReadStream struct {
	read_ *bufio.Reader
}

func NewReadStream(read io.Reader) *ReadStream {
	return &ReadStream{read_: bufio.NewReader(read)}
}

func (stream *ReadStream) ReadLine() ([]byte, error) {
	buffer := bytes.Buffer{}
	for {
		data, is, err := stream.read_.ReadLine()
		if err != nil {
			return data, err
		}
		if is {
			if len(data) > 0 {
				buffer.Write(data)
			}
		} else {
			buffer.Write(data)
			return buffer.Bytes(), nil
		}
	}
	return nil, nil
}
func (stream *ReadStream) read(len int) ([]byte, error) {
	data := make([]byte, len)
	var l = 0
	for l < len {
		n, err := stream.read_.Read(data[l:])
		if err != nil {
			return nil, err
		}
		l += n
	}
	return data, nil
}
func (stream *ReadStream) readUint(len uint32) ([]byte, error) {
	data := make([]byte, len)
	var l uint32 = 0
	for l < len {
		n, err := stream.read_.Read(data[l:])
		if err != nil {
			return nil, err
		}
		l += (uint32)(n)
	}
	return data, nil
}
func (stream *ReadStream) ReadUintBytes(len uint32) ([]byte, error) {
	return stream.readUint(len)
}

func (stream *ReadStream) ReadBytes(len int) ([]byte, error) {
	return stream.read(len)
}
func (stream *ReadStream) ReadByte() (byte, error) {
	return stream.read_.ReadByte()
}

type WriteStream struct {
	write_ *bufio.Writer
}

func NewWriteStream(write io.Writer) *WriteStream {
	return &WriteStream{write_: bufio.NewWriter(write)}
}

func (stream *WriteStream) Write(data []byte) (int, error) {
	return stream.write_.Write(data)
}
func (stream *WriteStream) Flush() error {
	return stream.write_.Flush()
}
