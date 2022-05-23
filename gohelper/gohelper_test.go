package gohelper

import (
	"testing"
)

func TestName(t *testing.T) {
	Run(func() {

		panic("+++++++")
	})
	Run(func() {
		panic("+++++++")

	})
	println("!!!!!!!!!!!!!")
	for{
		value:=Stack()
		println("===========================================")
		println(value)
	}
}
