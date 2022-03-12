package io

import (
	"bufio"
	"bytes"
	"io"
	"net"
)
func ReadAll(read io.Reader)([]byte, error){
	return io.ReadAll(read)
}
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
func (stream *ReadStream) ReadLineLimit(limit int) ([]byte, error){
	buffer := bytes.Buffer{}
	for {
		data, is, err := stream.read_.ReadLine()
		if err != nil {
			return data, err
		}
		if buffer.Len()+ len(data)>limit {
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
	return stream.read(len)
}
func (stream *ReadStream) ReadByte() (byte, error) {
	return stream.read_.ReadByte()
}

func (stream *ReadStream) Read(data []byte) (int,error) {
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
	isManualClose bool
}

func NewIOStream(cnn *net.TCPConn) *NetStream {
	var sm = &NetStream{conn: cnn, isManualClose: false}
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
func (stream *NetStream) WriteAndFlush(data []byte)  {
	stream.Write(data)
	stream.Flush()
}
func (stream *NetStream)ReadFunc(f func([]byte) bool ){
	data:=make([]byte,8192)
	go func() {
		for{
			num,err:=stream.ReadStream.Read(data)
			if err!=nil{
				break
			}else{
				if !f(data[0:num]){
					break
				}
			}
		}
	}()
}