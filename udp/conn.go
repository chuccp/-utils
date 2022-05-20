package udp

import (
	"bytes"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/queue"
	"github.com/chuccp/utils/udp/util"
	"github.com/chuccp/utils/udp/wire"
	"net"
	"sync"
	"time"
)

type connStore struct {
	connMap *sync.Map
}

type receivedPacket struct {
	data    []byte
	rcvTime time.Time
	header  *wire.Header
}

func (s *connStore) getGroup(conn *net.UDPConn, remoteAddr net.Addr) *GroupConn {
	var localAddr = conn.LocalAddr()
	key := "remote:" + remoteAddr.String() + "|" + "local:" + localAddr.String()
	v, ok := s.connMap.Load(key)
	if !ok {
		gp := NewGroupConn(remoteAddr, localAddr)
		s.connMap.Store(key, gp)

	}
	return v.(*GroupConn)
}

func (s *connStore) getConn(rawConn *net.UDPConn, remoteAddr net.Addr, isClient bool) (*QuicConn, error, bool) {
	var localAddr = rawConn.LocalAddr()
	key := "remote:" + remoteAddr.String() + "|" + "local:" + localAddr.String()
	v, ok := s.connMap.Load(key)
	if !ok {
		gp := NewGroupConn(remoteAddr, localAddr)
		s.connMap.Store(key, gp)
		if isClient {
			conn, err := gp.GetOrCreateClient(rawConn)
			return conn, err, true
		} else {
			conn := gp.GetOrCreateServer(rawConn)
			return conn, nil, true
		}
	} else {
		gp := v.(*GroupConn)
		if isClient {
			conn, err := gp.GetOrCreateClient(rawConn)
			return conn, err, true
		} else {
			conn := gp.GetOrCreateServer(rawConn)
			return conn, nil, true
		}
	}
	return nil, nil, false

}

func newConnStore() *connStore {
	return &connStore{connMap: new(sync.Map)}
}

type GroupConn struct {
	client     *QuicConn
	server     *QuicConn
	remoteAddr net.Addr
	localAddr  net.Addr
}

func NewGroupConn(remoteAddr net.Addr, localAddr net.Addr) *GroupConn {
	return &GroupConn{remoteAddr: remoteAddr, localAddr: localAddr}
}
func (group *GroupConn) Write(data []byte) (*QuicConn, bool) {

	group.handlePacket(data)

	return nil, false
}
func (group *GroupConn) handlePacket(data []byte) error {
	log.Info("读取数据 handlePacket ", wire.ConnectionID(data))
	header, err := wire.ParsePacket(data)
	if err != nil {
		return err
	}
	var pack = &receivedPacket{data: data, rcvTime: time.Now(), header: header}
	if header.IsLongHeader {
		switch header.Type {
		case wire.PacketTypeRetry:
			{
				client := group.client
				if client != nil {
					client.push(pack)
				}
			}
		case wire.PacketTypeInitial:
			{

			}

		}
	}
	return nil

}

func (group *GroupConn) GetOrCreateClient(conn *net.UDPConn) (*QuicConn, error) {
	if group.client == nil {
		group.client = newConn(group.remoteAddr, group.localAddr, true, conn)
		err := group.client.run()
		if err != nil {
			return nil, err
		}
	}
	return group.client, nil
}
func (group *GroupConn) GetOrCreateServer(rawConn *net.UDPConn) *QuicConn {

	if group.server == nil {
		group.server = newConn(group.remoteAddr, group.localAddr, false, rawConn)
	}

	return group.server
}
func (group *GroupConn) RemoteAddr() net.Addr {
	return group.remoteAddr
}
func (group *GroupConn) LocalAddr() net.Addr {
	return group.localAddr
}

type QuicConn struct {
	queue             *queue.VQueue
	buffer            *bytes.Buffer
	conn              *net.UDPConn
	IsClient          bool
	remoteAddr        net.Addr
	localAddr         net.Addr
	desConnectionID   wire.ConnectionID
	handShakeComplete bool
	handShakeProgress wire.HandShakeProgress
	handShakePacket   chan *receivedPacket
	PacketNum         wire.ByteCount
}

func (c *QuicConn) push(packet *receivedPacket) {

	if c.handShakeComplete {

	}

	c.queue.Offer(packet)
}

func (c *QuicConn) Read(p []byte) (n int, err error) {
	//alen := len(p)
	//var start = 0
	//if c.buffer.Len() > 0 {
	//	start, _ = c.buffer.Read(p)
	//	if start == alen {
	//		return start, err
	//	}
	//}
	//rDataLen := alen - start
	//data, _ := c.queue.Poll()
	//if data == nil {
	//	return 0, io.EOF
	//}
	//rData, _ := data.([]byte)
	//rLen := len(rData)
	//if rDataLen >= rLen {
	//	copy(p[start:], rData)
	//	return start + rLen, nil
	//}
	//copy(p[start:], rData[0:rDataLen])
	//c.buffer.Write(rData[rDataLen:])
	return 0, nil
}
func (c *QuicConn) Write(p []byte) (n int, err error) {
	return c.conn.WriteTo(p, c.remoteAddr)
}
func newConn(remoteAddr net.Addr, localAddr net.Addr, isClient bool, rawConn *net.UDPConn) *QuicConn {
	c := &QuicConn{queue: queue.NewVQueue(),
		remoteAddr:        remoteAddr,
		localAddr:         localAddr,
		buffer:            new(bytes.Buffer),
		conn:              rawConn,
		IsClient:          isClient,
		handShakeComplete: false}
	c.handShakePacket = make(chan *receivedPacket)
	c.PacketNum = 0
	return c
}

func (c *QuicConn) readBuffer() {
	func() {
		for {
			rev, _ := c.queue.Poll()
			rcp := rev.(*receivedPacket)
			if rcp.header.IsLongHeader {
				err := c.handleLongPacket(rcp)
				if err != nil {
					break
				}
			}
		}
	}()
}
func (c *QuicConn) handleLongPacket(rcp *receivedPacket) error {
	if !c.handShakeComplete {
		c.handShakePacket <- rcp
	}
	return nil
}
func (c *QuicConn) run() error {
	go c.readBuffer()
	return c.Initial()
}

func (c *QuicConn) Initial() error {
	log.Info("发送Initial")
	if c.IsClient {

		desConnectionID, err := wire.GenerateConnectionID(16)
		if err != nil {
			return err
		}
		initialParameter := &wire.InitialParameter{}
		initialParameter.ConnectionID = desConnectionID
		initialParameter.PacketType = wire.PacketTypeInitial
		initialParameter.PacketNum = c.PacketNum
		initialParameter.Token = []byte{}
		initialParameter.Random, _ = util.RandId(32)
		initialParameter.NextProtos = []string{"quic-echo-example", "cooge"}
		pack := wire.Initial(initialParameter)
		data := pack.Bytes()
		log.Info("=====", wire.ConnectionID(data))
		_, err1 := c.Write(data)
		c.handShakeProgress = wire.WaitRetry
		retryPacket := <-c.handShakePacket
		if retryPacket.header.IsLongHeader && retryPacket.header.Type == wire.PacketTypeRetry {
			if wire.HandleRetryPacket(retryPacket.data, initialParameter.ConnectionID) {
				c.PacketNum++
				token := retryPacket.data[retryPacket.header.ParsedLen : len(retryPacket.data)-16]
				log.Info("token:", wire.ConnectionID(token))
				initialParameter.Token = token
				initialParameter.PacketNum = c.PacketNum
				initialParameter.ConnectionID = retryPacket.header.SourceConnId
				pack := wire.Initial(initialParameter)
				data := pack.Bytes()
				_, err1 := c.Write(data)
				if err1 != nil {
					return err1
				}
			} else {
				log.Info("重试包验证失败")
				return util.RetryVerifyError
			}
		} else {
			return util.ProtocolError
		}
		return err1
	}
	return nil
}
