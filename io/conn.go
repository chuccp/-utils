package io

import (
	"github.com/chuccp/utils/log"
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
	log.InfoF("创建连接 {}",x.addr.String())
	conn, err := net.DialTCP("tcp", nil, x.addr)
	if err != nil {
		return nil,err
	}
	x.stream = NewIOStream(conn)
	return x.stream,nil
}
func (x *XConn) Close() {
	x.stream.Close()
}
func (x *XConn) LocalAddress() *net.TCPAddr {
	return x.stream.GetLocalAddress()
}
func (x *XConn) RemoteAddress() *net.TCPAddr {
	return x.stream.GetRemoteAddress()
}

func (x *XConn) WriteAndFlush(data []byte)  {
	x.stream.Write(data)
	x.stream.Flush()
}
func (x *XConn)Read(f func([]byte) bool ){
	data:=make([]byte,8192)
	go func() {
		for{
			num,err:=x.stream.Read(data)
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
