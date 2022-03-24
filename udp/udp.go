package udp

import (
	"io"
	"net"
	"strconv"
)

type Config struct {

}

type Listener struct {
	conn       *net.UDPConn
	connStore  *connStore
}

func (l *Listener) Accept() (*Conn, error) {
	for {
		data:=make([]byte,MaxPacketBufferSize)
		_,remoteAddr, err := l.conn.ReadFromUDP(data)
		if err != nil {
			return nil, err
		}
		conn, flag := l.connStore.getConn(l.conn,remoteAddr)
		conn.IsClient = false
		conn.push(data)
		if flag {
			return conn, nil
		}
	}
}
func (l *Listener) GetClientConn(address string) (*Conn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, flag :=l.connStore.getConn(l.conn,udpAddr)
	conn.IsClient = true
	if flag{
		return conn, nil
	}
	return nil, io.EOF
}

func newListener(conn *net.UDPConn) *Listener {
	return &Listener{conn: conn, connStore: newConnStore()}
}
func ListenAddr(port int) (*Listener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return newListener(conn), nil
}
