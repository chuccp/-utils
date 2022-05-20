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
	log.Print("新建协程",atomic.AddInt32(&h.coroutineNum,1))

}
func (h *GoHelper)end(){
	log.Print("协程运行结束",atomic.AddInt32(&h.coroutineNum,-1))
}
func (h *GoHelper)  panic() {
	log.Print("协程运行异常",atomic.AddInt32(&h.coroutineNum,-1))
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