package math

import (
	"math/rand"
	"time"
)

func RandInt() uint32 {
	num := rand.Intn(1024)
	return millisecond()<<10 | (uint32(num))
}
func millisecond() uint32 {
	ms := time.Now().UnixNano() / 1e6
	return uint32(ms)
}
