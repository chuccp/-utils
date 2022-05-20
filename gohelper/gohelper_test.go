package gohelper

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	Run(func() {

		panic("+++++++")
	})
	Run(func() {
		time.Sleep(time.Second)

	})
	println("!!!!!!!!!!!!!")
	for{
		value:=Stack()
		println("===========================================")
		println(value)
	}
}
