package udp

import (
	"github.com/chuccp/utils/log"
	"math/rand"
	"strconv"
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
					var data = make([]byte, 1024)
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
	time.Sleep(time.Second * 5)
	for {
		conn, err := listen.GetClientConn("127.0.0.1:8090")
		if err == nil {
			conn.Write([]byte(strconv.Itoa(rand.Int())+"==="+strconv.Itoa(rand.Int())))
			time.Sleep(time.Second*2)
		} else {
			log.Info(err)
			break
		}
	}

}
