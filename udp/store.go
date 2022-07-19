package udp

import (
	"net"
	"sync"
)

type Store struct {
	connMap  *sync.Map
}

func (s *Store) Load(addr net.Addr)(*RawConn,bool)  {
	 k:=addr.String()
	v,ok:=s.connMap.Load(k)
	if ok{
		return v.(*RawConn),true
	}else{
		cc:= newRawConn()
		actual,loaded:=s.connMap.LoadOrStore(k,cc)
		if loaded{
			return  actual.(*RawConn),true
		}else{
			return  v.(*RawConn),false
		}
	}
}
func NewStore() *Store {
	return &Store{connMap: new(sync.Map)}
}