package log

import (
	"bufio"
	"bytes"
	"log"
	"runtime/debug"
	"testing"
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
func run() {
	log.Panic("ooooo")
}
func recordRecover() {
	err := recover()
	if err != nil {
		var err = debug.Stack()
		var buffer = bytes.NewBuffer(err)
		var br = bufio.NewReader(buffer)
		for {
			line, fa, errLine := br.ReadLine()
			if errLine != nil {
				break
			}
			if !fa {
				v := string(line)
				log.Println("系统级错误  ", v)

			} else {
				break
			}
		}
	}
}
func TestInfo3(t *testing.T) {
	config := GetDefaultConfig()
	config.SetLevel(DebugLevel)
	config.AddFileConfig("log/${time:2006-01-02-15-04}-${line:2000}-${size:200mb}-${level}.log", ErrorLevel)

	ErrorF("=============")
}
