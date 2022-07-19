package util

import "net"

func CreateUdpAddr(address string) (net.Addr,error) {
	return net.ResolveUDPAddr("udp",address)
}