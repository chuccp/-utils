package udp

import (
	"bytes"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/queue"
	"io"
	"net"
	"sync"
)

type connStore struct {
	connMap *sync.Map
}

func (s *connStore) getConn(conn *net.UDPConn, remoteAddr net.Addr, isClient bool) (*Conn, bool) {
	var localAddr = conn.LocalAddr()
	key := "remote:" + remoteAddr.String() + "|" + "local:" + localAddr.String()
	v, ok := s.connMap.Load(key)
	if !ok {
		gp := NewGroupConn(remoteAddr, localAddr)
		s.connMap.Store(key, gp)
		if isClient {
			return gp.GetOrCreateClient(conn), true
		} else {
			return gp.GetOrCreateClient(conn), true
		}
	}else{
		gp := v.(*GroupConn)
		if isClient {
			return gp.GetOrCreateClient(conn), true
		} else {
			return gp.GetOrCreateClient(conn), true
		}
	}
	return nil, false

}

func newConnStore() *connStore {
	return &connStore{connMap: new(sync.Map)}
}

type GroupConn struct {
	client     *Conn
	server     *Conn
	remoteAddr net.Addr
	localAddr  net.Addr
}

func NewGroupConn(remoteAddr net.Addr, localAddr net.Addr) *GroupConn {
	return &GroupConn{remoteAddr: remoteAddr, localAddr: localAddr}
}
func (group *GroupConn) GetOrCreateClient(conn *net.UDPConn) *Conn {
	if group.client == nil {
		group.client = newConn(group.remoteAddr, group.localAddr, true, conn)
		group.client.Initial()
	}
	return group.client
}
func (group *GroupConn) GetOrCreateServer(conn *net.UDPConn) *Conn {

	if group.server == nil {
		group.server = newConn(group.remoteAddr, group.localAddr, false, conn)
	}

	return group.server
}
func (group *GroupConn) RemoteAddr() net.Addr {
	return group.remoteAddr
}
func (group *GroupConn) LocalAddr() net.Addr {
	return group.localAddr
}

type Conn struct {
	queue      *queue.VQueue
	buffer     *bytes.Buffer
	conn       *net.UDPConn
	IsClient   bool
	remoteAddr net.Addr
	localAddr  net.Addr
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
	log.Info("写数据：",p," ",c.remoteAddr)
	return c.conn.WriteTo(p, c.remoteAddr)
}
func newConn(remoteAddr net.Addr, localAddr net.Addr, isClient bool, conn *net.UDPConn) *Conn {
	return &Conn{queue: queue.NewVQueue(), remoteAddr: remoteAddr, localAddr: localAddr, buffer: new(bytes.Buffer), conn: conn, IsClient: isClient}
}

func (c *Conn) Initial() error {
	if c.IsClient{
		//log.Info("!!!写数据")
		//data:=make([]byte,MaxPacketBufferSize)
		//f,_:=file.NewFile("C:\\Users\\cooge\\Documents\\quic\\Initial4.bin")
		//n,_:=f.ReadBytes(data)
		////connId, err := GenerateConnectionID(16)
		////if err != nil {
		////	return err
		////}
		////pack := Initial(connId)
		////data := pack.Bytes()
		//num,err := c.Write(data[0:n])
		//log.Info(num," ",err)

		connId:=[]byte{0x8c,0x07,0x0c,0x38,0x6e,0xf8,0xef,0x21,0x9c,0x6d,0xba,0x73,0xed,0xa7,0x4f}
		pack:=Initial(connId)
		data:=pack.Bytes()
		log.Info(ConnectionID(data))
		_,err := c.Write(data)
		return err
	}
	return nil
}
