package udp

import (
	"github.com/chuccp/utils/log"
	"testing"
)

func TestName(t *testing.T) {
	listen, err := ListenAddr(8090)
	if err != nil {
		log.Info(err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				break
			}
			go func() {
				for {
					var data = make([]byte, MaxPacketBufferSize)
					num, err := conn.Read(data)
					if err == nil {
						log.Info(num, string(data[0:num]))
					} else {
						log.Info(err)
					}
				}
			}()
		}
	}()
	log.Info("获取客户端")
	for {
		_, err := listen.GetClientConn("129.211.17.31:8086")
		if err == nil {
		} else {
			log.Info(err)
			break
		}
	}

}
