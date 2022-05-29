package io

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"
)

func ReadAll(read io.Reader) ([]byte, error) {
	return io.ReadAll(read)
}

type ReadStream struct {
	read_  *bufio.Reader
	offset uint16
}

func NewReadStream(read io.Reader) *ReadStream {
	return &ReadStream{read_: bufio.NewReader(read)}
}
func NewReadBytesStream(data []byte) *ReadStream {
	return NewReadStream(bytes.NewReader(data))
}
func (stream *ReadStream) Offset() uint16 {
	return stream.offset
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
func (stream *ReadStream) ReadLineLimit(limit int) ([]byte, error) {
	buffer := bytes.Buffer{}
	for {
		data, is, err := stream.read_.ReadLine()
		if err != nil {
			return data, err
		}
		if buffer.Len()+len(data) > limit {
			return nil, bufio.ErrBufferFull
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
	stream.offset =  stream.offset+uint16(len)
	return stream.read(len)
}
func (stream *ReadStream) Read2Uint16() (uint16, error) {
	data, err := stream.read(2)
	if err != nil {
		return 0, err
	} else {
		stream.offset = stream.offset + 2
		return uint16(data[0])<<8 | uint16(data[1]), nil
	}
}
func (stream *ReadStream) Read3Uint32() (uint32, error) {
	data, err := stream.read(3)
	if err != nil {
		return 0, err
	} else {
		stream.offset = stream.offset + 3
		return uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]), nil
	}
}
func (stream *ReadStream) Read4Uint32() (uint32, error) {
	data, err := stream.read(4)
	if err != nil {
		return 0, err
	} else {
		stream.offset = stream.offset + 4
		return uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]), nil
	}
}
func (stream *ReadStream) ReadByte() (byte, error) {
	stream.offset++
	return stream.read_.ReadByte()
}
func (stream *ReadStream) ReadUint8() (uint8, error) {
	stream.offset++
	return stream.read_.ReadByte()
}

func (stream *ReadStream) Read(data []byte) (int, error) {
	return stream.read_.Read(data)
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

type NetStream struct {
	conn *net.TCPConn
	*ReadStream
	*WriteStream
	once          *sync.Once
	isManualClose bool
}

func NewStream(cnn *net.TCPConn) *NetStream {
	var sm = &NetStream{conn: cnn, isManualClose: false, once: new(sync.Once)}
	sm.WriteStream = NewWriteStream(cnn)
	sm.ReadStream = NewReadStream(cnn)
	return sm
}
func (stream *NetStream) GetLocalAddress() *net.TCPAddr {
	if stream.conn.LocalAddr() == nil {
		return nil
	}
	return stream.conn.LocalAddr().(*net.TCPAddr)
}
func (stream *NetStream) GetRemoteAddress() *net.TCPAddr {
	return stream.conn.RemoteAddr().(*net.TCPAddr)
}

func (stream *NetStream) ManualClose() {
	stream.isManualClose = true
	stream.conn.Close()
}
func (stream *NetStream) IsManualClose() bool {
	return stream.isManualClose
}
func (stream *NetStream) WriteAndFlush(data []byte) (num int, err error) {
	num, err = stream.Write(data)
	if err != nil {
		return
	}
	err = stream.Flush()
	return
}
func (stream *NetStream) ReadFunc(f func([]byte) bool, close func()) {
	data := make([]byte, 8192)
	go func() {
		for {
			num, err := stream.ReadStream.Read(data)
			if err != nil {
				break
			} else {
				if !f(data[0:num]) {
					break
				}
			}
		}
		stream.once.Do(func() {
			close()
		})
	}()
}
