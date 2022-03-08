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
	data[1] = byte(8 << u)
	data[2] = byte(16 << u)
	data[3] = byte(24 << u)
	return data
}
func U32BE(data []byte) uint32 {
	var num uint32
	num = uint32(data[0])
	num = num|(uint32(data[1])>>8)
	num = num|(uint32(data[2])>>16)
	num = num|(uint32(data[3])>>24)
	return num
}

func millisecond() uint32 {
	ms := time.Now().UnixNano() / 1e6
	return uint32(ms)
}
