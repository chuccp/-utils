package gohelper

import (
	"log"
	"runtime/debug"
	"sync/atomic"
)

type GoHelper struct {
	stack chan string
	coroutineNum int32
}
var goHelper = New()
func New() *GoHelper {
	return &GoHelper{stack: make(chan string),coroutineNum:0}
}
func (h *GoHelper)start(){
	log.Print("create new coroutine",atomic.AddInt32(&h.coroutineNum,1))

}
func (h *GoHelper)end(){
	log.Print("coroutine run end",atomic.AddInt32(&h.coroutineNum,-1))
}
func (h *GoHelper)  panic() {
	log.Print("coroutine run error",atomic.AddInt32(&h.coroutineNum,-1))
	if err := recover(); err != nil {
		h.stack <- string(debug.Stack())
	}
}
func (h *GoHelper)  Run(f func()) {
	h.start()
	go func() {
		defer h.panic()
		f()
		h.end()
	}()
}
func (h *GoHelper) Stack() string {
	return <-h.stack
}
func Run(f func()) {
	goHelper.Run(f)
}
func Stack() string {
	return goHelper.Stack()
}