package udp

import (
	"net"
	"strconv"
)

type Listener struct {
	conn       *net.UDPConn
	readStream *ReadStream
	connStore  *connStore
}

func (l *Listener) Accept() (*Conn, error) {

	for {
		data,remoteAddr, err := l.readStream.Read(1024)
		if err != nil {
			return nil, err
		}
		conn, flag := l.connStore.getConn(l.conn,remoteAddr)
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
	return newConn(udpAddr, l.conn.LocalAddr(), l.conn), nil
}

func newListener(conn *net.UDPConn) *Listener {
	return &Listener{conn: conn, readStream: newReadStream(conn), connStore: newConnStore()}
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
