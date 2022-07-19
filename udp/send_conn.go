package udp

import "net"

type rawPacket struct {
	remoteAddr net.Addr
	localAddr net.Addr
	oob        []byte
}

func NewRawPacket(remoteAddr string,localAddr string,data []byte)*rawPacket  {
	rAddr, _ := net.ResolveUDPAddr("udp", remoteAddr)
	lAddr,_:=net.ResolveUDPAddr("udp",localAddr)
	return &rawPacket{remoteAddr:rAddr,localAddr:lAddr,oob: data}
}