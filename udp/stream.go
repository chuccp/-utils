package udp

import (
	"net"
)

type ReadStream struct {
	conn  *net.UDPConn
}
func newReadStream(conn  *net.UDPConn)*ReadStream{
	return &ReadStream{conn:conn}
}

func (readStream *ReadStream) Read(num int)([]byte,*net.UDPAddr,error)  {
	data:=make([]byte,num)
	num,addr,err:=readStream.conn.ReadFromUDP(data)
	return data[0:num],addr, err

}