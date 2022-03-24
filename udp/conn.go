package udp

import (
	"bytes"
	"github.com/chuccp/utils/queue"
	"io"
	"net"
	"sync"
)

type connStore struct {
	connMap *sync.Map
}

func (s *connStore) getConn(conn *net.UDPConn,remoteAddr net.Addr) (*Conn, bool) {
	var localAddr = conn.LocalAddr()
	if remoteAddr.String()==localAddr.String(){
		panic("remoteAddr can't same with localAddr")
	}
	key := "remote:" + remoteAddr.String() + "|" + "local:" + localAddr.String()
	v, ok := s.connMap.Load(key)
	if ok {
		return v.(*Conn), false
	} else {
		c := newConn(remoteAddr, localAddr, conn)
		s.connMap.Store(key, c)
		return c, true
	}
}

func newConnStore() *connStore {
	return &connStore{connMap: new(sync.Map)}
}

type Conn struct {
	queue      *queue.VQueue
	remoteAddr net.Addr
	localAddr  net.Addr
	buffer     *bytes.Buffer
	conn       *net.UDPConn
	IsClient bool
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}
func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *Conn) push(data []byte) {
	c.queue.Offer(data)
}

func (c *Conn) Read(p []byte) (n int, err error) {
	alen := len(p)
	var start = 0
	if c.buffer.Len() > 0 {
		start, _ = c.buffer.Read(p)
		if start == alen {
			return start, err
		}
	}
	rDataLen := alen - start
	data, _ := c.queue.Poll()
	if data == nil {
		return 0, io.EOF
	}
	rData, _ := data.([]byte)
	rLen := len(rData)
	if rDataLen >= rLen {
		copy(p[start:], rData)
		return start + rLen, nil
	}
	copy(p[start:], rData[0:rDataLen])
	c.buffer.Write(rData[rDataLen:])
	return alen, nil
}
func (c *Conn) Write(p []byte) (n int, err error) {
	return c.conn.WriteTo(p, c.remoteAddr)
}
func newConn(remoteAddr net.Addr, localAddr net.Addr, conn *net.UDPConn) *Conn {
	return &Conn{queue: queue.NewVQueue(), remoteAddr: remoteAddr, localAddr: localAddr, buffer: new(bytes.Buffer), conn: conn}
}
