package udp

import (
	"net"
)



type baseServer struct {
	udpConn *net.UDPConn
	store *Store
}

func (s *baseServer) listen(port int) error {

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return err
	}
	s.udpConn = udpConn
	return nil
}
func (s *baseServer)Accept()(*RawConn,error){
	for{
		data:=make([]byte, MaxPacketBufferSize)
		num, remoteLocal, err := s.udpConn.ReadFromUDP(data)
		if err != nil {
			return nil, err
		}else{
			rc,flag:=s.store.Load(remoteLocal)
			if flag{
				rc.push(data[0:num])
				return rc,nil
			}else{
				rc.push(data[0:num])
			}
		}
	}
	return nil,nil
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
