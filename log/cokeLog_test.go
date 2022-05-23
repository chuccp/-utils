package log

import (
	"bufio"
	"log"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	config := GetDefaultConfig()
	var b = []byte{1, 2, 3, 4}
	config.SetLevel(DebugLevel)
	Debug(b, bufio.ErrInvalidUnreadByte)
	Debug(b, bufio.ErrInvalidUnreadByte)
	Debug(b, bufio.ErrInvalidUnreadByte)

}
func TestInfo2(t *testing.T) {
	config := GetDefaultConfig()
	var b = []byte{1, 2, 3, 4}
	config.SetLevel(DebugLevel)
	log.Println(b, bufio.ErrInvalidUnreadByte)
	log.Println(b, bufio.ErrInvalidUnreadByte)
	log.Println(b, bufio.ErrInvalidUnreadByte)

}

func TestInfo3(t *testing.T) {
	config := GetDefaultConfig()
	config.SetLevel(DebugLevel)
	config.AddFileConfig("log/${time:2006-01-02-15-04}-${line:2000}-${size:200mb}-${level}.log", DebugLevel)

	Error("=============")
	Info("=============")

	time.Sleep(time.Hour)

}
