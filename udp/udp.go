package udp

import (
	"github.com/chuccp/utils/log"
	"io"
	"net"
	"strconv"
)

type Config struct {
}

type Listener struct {
	conn      *net.UDPConn
	connStore *connStore
	chanConn  chan *Conn
}

func (l *Listener) readUDP() {
	log.Info("开始读取信息")
	for {
		data := make([]byte, MaxPacketBufferSize)
		num, remoteAddr, err := l.conn.ReadFromUDP(data)
		log.Info("读取数据",num,remoteAddr,err)
		if err == nil {
			groupConn := l.connStore.getGroup(l.conn, remoteAddr)
			conn, flag := groupConn.Write(data[:num])
			if flag {
				l.chanConn <- conn
			}
		} else {
			break
		}
	}
	l.chanConn <- nil
}

func (l *Listener) Accept() (*Conn, error) {
	c := <-l.chanConn
	if c == nil {
		return nil, io.EOF
	}
	return c, nil
}
func (l *Listener) GetClientConn(address string) (*Conn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err, flag := l.connStore.getConn(l.conn, udpAddr, true)
	if err != nil {
		return nil, err
	}
	conn.IsClient = true
	if flag {
		return conn, nil
	}
	return nil, io.EOF
}

func newListener(conn *net.UDPConn) *Listener {
	chanConn := make(chan *Conn)
	listener:= &Listener{conn: conn, connStore: newConnStore(), chanConn: chanConn}
	go listener.readUDP()
	return listener
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
