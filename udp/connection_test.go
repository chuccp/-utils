package udp

import "testing"

func TestName(t *testing.T) {

	server, err := listen(8685)
	if err != nil {
		return
	}
	for{
		conn, err := server.Accept()
		if err != nil {
			return
		}
		go
			func() {
				for{
					var data =make([]byte, MaxPacketBufferSize)
					_, err := conn.Read(data)
					if err != nil {
						break
					}
				}
			}()
	}
}
