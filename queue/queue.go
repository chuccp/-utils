package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Queue struct {
	input   *element
	output  *element
	ch      chan bool
	waitNum int32
	num     int32
	lock    *sync.RWMutex
	rLock   *sync.Mutex
	timer   *timer
}

func NewQueue() *Queue {
	return &Queue{ch: make(chan bool), waitNum: 0, num: 0, lock: new(sync.RWMutex), rLock: new(sync.Mutex)}
}
func (queue *Queue) Offer(value interface{}) (num int32) {
	ele := getElement(value)
	queue.lock.Lock()
	if atomic.CompareAndSwapInt32(&queue.num, 0, 1) {
		queue.input = ele
		queue.output = ele
		num = 1
	} else {
		queue.input.next = ele
		queue.input = ele
		num = atomic.AddInt32(&queue.num, 1)
	}
	if queue.waitNum > 0 {
		atomic.AddInt32(&queue.waitNum, -1)
		queue.lock.Unlock()
		queue.ch <- true
	} else {
		queue.lock.Unlock()
	}
	return
}
func (queue *Queue) Num() int32 {
	return queue.num
}
func (queue *Queue) Poll() (value interface{}, num int32) {
	for {
		queue.lock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				value, num = queue.readOne()
				queue.lock.Unlock()
				return
			} else {
				queue.lock.Unlock()
				queue.rLock.Lock()
				val, n, last := queue.readGtOne()
				if last {
					queue.rLock.Unlock()
				} else {
					queue.rLock.Unlock()
					return val, n
				}
			}
		} else {
			queue.waitNum++
			queue.lock.Unlock()
			<-queue.ch
		}
	}
}
func (queue *Queue) Peek() (value interface{}, num int32) {
	queue.lock.RLock()
	num = queue.num
	if queue.num > 0 {
		value = queue.output.value
		queue.lock.RUnlock()
		return
	} else {
		queue.lock.RUnlock()
		return nil, 0
	}
}
func (queue *Queue) readOne() (value interface{}, num int32) {
	var ele = queue.output
	value = ele.value
	num = atomic.AddInt32(&queue.num, -1)
	freeElement(ele)
	return value, num
}
func (queue *Queue) readGtOne() (value interface{}, num int32, isLast bool) {
	var ele = queue.output
	if ele.next == nil {
		return nil, 0, true
	}
	value = ele.value
	queue.output = ele.next
	num = atomic.AddInt32(&queue.num, -1)
	freeElement(ele)
	return value, num, false
}

func (queue *Queue) Dequeue(ctx context.Context) (value interface{}, num int32, colse bool) {

	for {
		queue.lock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				value, num = queue.readOne()
				queue.lock.Unlock()
				return
			} else {
				queue.lock.Unlock()
				queue.rLock.Lock()
				val, n, last := queue.readGtOne()
				if last {
					queue.rLock.Unlock()
				} else {
					queue.rLock.Unlock()
					return val, n, false
				}
			}
		} else {
			queue.waitNum++
			queue.lock.Unlock()
			var op = getOperate(ctx)
			go func() {
				fa := op.wait()
				queue.lock.Lock()
				if !fa && queue.waitNum > 0 {
					queue.waitNum--
					queue.lock.Unlock()
					queue.ch <- !fa
				} else {
					queue.lock.Unlock()
				}
			}()
			flag := <-queue.ch
			freeOperate(op)
			if !flag {
				return nil, 0, true
			}
		}
	}
}

func (queue *Queue) Take(duration time.Duration) (value interface{}, num int32) {
	for {
		queue.lock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				value, num = queue.readOne()
				queue.lock.Unlock()
				return
			} else {
				queue.lock.Unlock()
				queue.rLock.Lock()
				val, n, last := queue.readGtOne()
				if last {
					queue.rLock.Unlock()
				} else {
					queue.rLock.Unlock()
					return val, n
				}
			}
		} else {
			queue.waitNum++
			queue.lock.Unlock()
			tm := getTimer(duration)
			go func() {
				fa := tm.wait()
				if !fa {
					queue.lock.Lock()
					if queue.waitNum > 0 {
						queue.waitNum--
						queue.lock.Unlock()
						queue.ch <- false
					} else {
						queue.lock.Unlock()
					}
				}
			}()
			flag := <-queue.ch
			tm.end()
			freeTimer(tm)
			if !flag {
				return nil, 0
			}
		}
	}
}
