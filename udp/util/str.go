package util

import (
	"encoding/hex"
	"strings"
)

func SplitHexString(s string) (slice []byte) {
	for _, ss := range strings.Split(s, " ") {
		if ss[0:2] == "0x" {
			ss = ss[2:]
		}
		d, _ := hex.DecodeString(ss)
		slice = append(slice, d...)
	}
	return
}
