package udp

import (
	"net"
	"sync"
)

type Store struct {
	connMap  *sync.Map
}

func (s *Store) Load(addr net.Addr)(*ReceiveConn,bool)  {
	 k:=addr.String()
	v,ok:=s.connMap.Load(k)
	if ok{
		return v.(*ReceiveConn),true
	}else{
		cc:= newReceiveConn()
		actual,loaded:=s.connMap.LoadOrStore(k,cc)
		if loaded{
			return  actual.(*ReceiveConn),true
		}else{
			return  v.(*ReceiveConn),false
		}
	}
}
func NewStore() *Store {
	return &Store{connMap: new(sync.Map)}
}