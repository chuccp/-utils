package io

import (
	"net"
	"strconv"
)

type XConn struct {
	port   int
	host   string
	address string
	addr   *net.TCPAddr
	stream *NetStream
}

func NewXConn(host string, port int) *XConn {
	addr:= host+":"+strconv.Itoa(port)
	return NewXConn2(addr)
}
func NewXConn2(address string) *XConn {
	addr, _ := net.ResolveTCPAddr("tcp", address)
	return &XConn{port: addr.Port, host: addr.Network(), addr: addr}
}
func (x *XConn) Create() (*NetStream,error) {
	conn, err := net.DialTCP("tcp", nil, x.addr)
	if err != nil {
		return nil,err
	}
	x.stream = NewStream(conn)
	return x.stream,nil
}
func (x *XConn) Close() {
	x.stream.conn.Close()
}
func (x *XConn) LocalAddress() *net.TCPAddr {
	return x.stream.GetLocalAddress()
}
func (x *XConn) RemoteAddress() *net.TCPAddr {
	return x.stream.GetRemoteAddress()
}
