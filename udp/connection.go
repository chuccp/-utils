package udp

import (
	"github.com/chuccp/utils/udp/wire"
	"net"
	"time"
)

type baseServer struct {
	udpConn *net.UDPConn
	store   *Store
}

func (s *baseServer) listen(port int) error {

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return err
	}
	s.udpConn = udpConn
	return nil
}
func (s *baseServer) Accept() (*ReceiveConn, error) {
	for {
		data := make([]byte, MaxPacketBufferSize)
		num, remoteLocal, err := s.udpConn.ReadFromUDP(data)
		if err != nil {
			return nil, err
		} else {
			packet, b, err := s.parsePacket(data[0:num], remoteLocal)
			if err != nil {
				return nil, err
			}
			if b {
				return packet, nil
			}
		}
	}
	return nil, nil
}
func (s *baseServer) parsePacket(data []byte, remote *net.UDPAddr) (*ReceiveConn, bool, error) {
	currentTime := time.Now()
	var header wire.Header
	err := ParseHeader(data, &header)
	if err != nil {
		return nil, false, err
	}
	if header.IsLongHeader {
		if header.PacketType == wire.PacketTypeInitial {
			receivePacket := NewReceivePacket(&currentTime, data, &header)
			rc, flag := s.store.Load(remote)
			if flag {
				rc.push(receivePacket)
				return rc, flag, nil
			} else {
				rc.push(receivePacket)
				return nil, flag, nil
			}
		} else {
			return nil, false, err
		}
	} else {
		return nil, false, err
	}
}

func newBaseServer() *baseServer {
	return &baseServer{store: NewStore()}
}

func listen(port int) (*baseServer, error) {
	s := newBaseServer()
	err := s.listen(port)
	if err != nil {
		return nil, err
	}
	return s, nil
}
