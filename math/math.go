package math

import (
	"math/rand"
	"time"
)

func RandUInt32() uint32 {
	num := rand.Intn(1024)
	return millisecond()<<10 | (uint32(num))
}

func BEU32(u uint32) []byte {
	var data = []byte{0, 0, 0, 0}
	data[0] = byte(u)
	data[1] = byte(u >> 8)
	data[2] = byte(u >> 16)
	data[3] = byte(u >> 24)
	return data
}
func U32BE(data []byte) uint32 {
	var num uint32
	num = uint32(data[0])
	num = num | (uint32(data[1]) << 8)
	num = num | (uint32(data[2]) << 16)
	num = num | (uint32(data[3]) << 24)
	return num
}

func U32BE0To4(data []byte,len uint8) uint32 {
	switch len {
	case 1:
		return uint32(data[0])
	case 2:
		return uint32(data[0])<<8|uint32(data[1])
	case 3:
		return uint32(data[0])<<16|uint32(data[1])<<8|uint32(data[2])
	case 4:
		return uint32(data[0])<<24|uint32(data[1])<<16|uint32(data[2])<<8|uint32(data[3])
	}
	return 0
}


func millisecond() uint32 {
	ms := time.Now().UnixNano() / 1e6
	return uint32(ms)
}
