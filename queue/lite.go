package queue

import "sync/atomic"

type LiteQueue struct {
  values	chan interface{}
  num int32
}

func NewLiteQueue(num int) *LiteQueue {
	values:=make(chan interface{},num)
	return  &LiteQueue{values:values}
}
func (lite *LiteQueue)Offer(value interface{})int32  {
	nn:=atomic.AddInt32(&lite.num, 1)
	lite.values<-value
	return nn
}
func (lite *LiteQueue) poll()  (value interface{}, nu int32){
	v:=<-lite.values
	return v,atomic.AddInt32(&lite.num, -1)
}