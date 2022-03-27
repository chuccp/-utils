package udp

import (
	"github.com/chuccp/utils/log"
	"testing"
	"time"
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
					log.Info("#######:",data)
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
	time.Sleep(time.Second * 1)
	for {
		_, err := listen.GetClientConn("129.211.17.31:8086")
		if err == nil {
			//conn.Write([]byte(strconv.Itoa(rand.Int())+"==="+strconv.Itoa(rand.Int())))
			time.Sleep(time.Hour*2)
		} else {
			log.Info(err)
			break
		}
	}

}
