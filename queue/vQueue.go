package queue

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type ele struct {
	next   *ele
	fq     []interface{}
	wIndex int32
	rIndex int32
	cap    int32
	lCap   int32
}

func (q *ele) insert(value interface{}) int32 {
	num := atomic.AddInt32(&q.wIndex, 1)
	if num > q.lCap {
		return num
	}
	q.fq[num] = value
	return num
}
func (q *ele) get() (interface{}, int32) {
	num := atomic.AddInt32(&q.rIndex, 1)
	if num > q.lCap {
		return nil, num
	}
	for {
		n := q.fq[num]
		if n != nil {
			q.fq[num] = nil
			return n, num
		} else {
			runtime.Gosched()
		}
	}
}

func (q *ele) read(num int32) interface{} {
	return q.fq[num]
}

var poolEle = &sync.Pool{
	New: func() interface{} {
		return new(ele)
	},
}

func newEle(ca int32) *ele {
	ele := poolEle.Get().(*ele)
	ele.fq = make([]interface{}, ca)
	ele.cap = ca
	ele.wIndex = -1
	ele.rIndex = -1
	ele.lCap = ca - 1
	return ele
}
func freeEle(ele *ele) {
	ele.next = nil
	poolEle.Put(ele)
}

type timer struct {
	t     *time.Timer
	isEnd chan bool
}

func newTimer() *timer {
	return &timer{t: time.NewTimer(time.Second * 10), isEnd: make(chan bool)}
}
func (timer *timer) wait() bool {
	select {
	case <-timer.t.C:
		{
			return false
		}
	case fa := <-timer.isEnd:
		{
			return fa
		}
	}
	return false
}
func (timer *timer) end() {
	fa := timer.t.Stop()
	if fa {
		timer.isEnd <- true
	}
}
func (timer *timer) reset(duration time.Duration) {
	timer.t.Reset(duration)
}

var poolTimer = &sync.Pool{
	New: func() interface{} {
		return newTimer()
	},
}

func getTimer(duration time.Duration) *timer {
	ti := poolTimer.Get().(*timer)
	ti.reset(duration)
	return ti
}
func freeTimer(timer *timer) {
	poolTimer.Put(timer)
}

type operate struct {
	ctx   context.Context
}

func newOperate(ctx context.Context) *operate {
	return &operate{ctx: ctx}
}
func (op *operate) wait() bool {
	select {
	case <-op.ctx.Done():
		return true
	}
	return false
}
func (op *operate) isClose() bool {
	return op.ctx.Err() != nil
}

var poolOperate = &sync.Pool{
	New: func() interface{} {
		return &operate{}
	},
}

func getOperate(ctx context.Context) *operate {
	v, _ := poolOperate.Get().(*operate)
	v.ctx = ctx
	return v
}
func freeOperate(op *operate) {
	poolOperate.Put(op)
}

type VQueue struct {
	write   *ele
	read    *ele
	num     int32
	waitNum int32
	flag    chan bool
	cap     int32
	lCap    int32
	lock    *sync.Mutex
}

func (queue *VQueue) Offer(value interface{}) (nu int32) {

	for {
		num := queue.write.insert(value)
		if num < queue.cap {
			nu = atomic.AddInt32(&queue.num, 1)
			queue.lock.Lock()
			if atomic.LoadInt32(&queue.waitNum) > 0 {
				atomic.AddInt32(&queue.waitNum, -1)
				queue.lock.Unlock()
				queue.flag <- true
			} else {
				queue.lock.Unlock()
			}
			return
		} else {
			if num == queue.cap {
				queue.write.next = newEle(queue.cap)
				queue.write = queue.write.next
			} else {
				runtime.Gosched()
			}
		}
	}
}
func (queue *VQueue) poll() (value interface{}, nu int32, hasValue bool) {
	v, num := queue.read.get()
	if num < queue.cap {
		nu = atomic.AddInt32(&queue.num, -1)
		return v, nu, true
	} else {
		if num == queue.cap {
			r := queue.read
			for {
				if queue.read.next != nil {
					queue.read = queue.read.next
					break
				} else {
					runtime.Gosched()
				}
			}
			freeEle(r)
		} else {
			runtime.Gosched()
		}
	}
	return
}

func (queue *VQueue) Poll() (value interface{}, nu int32) {
	for {
		queue.lock.Lock()
		if atomic.LoadInt32(&queue.num) == 0 {
			atomic.AddInt32(&queue.waitNum, 1)
			queue.lock.Unlock()
			<-queue.flag
		} else {
			queue.lock.Unlock()
			v, num, has := queue.poll()
			if has {
				return v, num
			}
		}
	}
}

func (queue *VQueue) Dequeue(ctx context.Context) (value interface{}, num int32, hasClose bool) {
	for {
		queue.lock.Lock()
		if atomic.LoadInt32(&queue.num) == 0 {
			atomic.AddInt32(&queue.waitNum, 1)
			queue.lock.Unlock()
			var op = getOperate(ctx)
			go func() {
				fa := op.wait()
				queue.lock.Lock()
				if atomic.LoadInt32(&queue.waitNum) > 0 {
					atomic.AddInt32(&queue.waitNum, -1)
					queue.lock.Unlock()
					queue.flag <- !fa
				} else {
					queue.lock.Unlock()
				}
			}()
			flag := <-queue.flag
			freeOperate(op)
			if !flag {
				return nil, 0, true
			}
		} else {
			queue.lock.Unlock()
			v, num, has := queue.poll()
			if has {
				return v, num, false
			}
		}
	}
}
func (queue *VQueue) Num()int32{
	return atomic.LoadInt32(&queue.num)
}
func (queue *VQueue) Take(duration time.Duration) (value interface{}, num int32) {
	for {
		queue.lock.Lock()
		if atomic.LoadInt32(&queue.num) == 0 {
			atomic.AddInt32(&queue.waitNum, 1)
			queue.lock.Unlock()
			tm := getTimer(duration)
			go func() {
				fa := tm.wait()
				if !fa {
					queue.lock.Lock()
					if atomic.LoadInt32(&queue.waitNum) > 0 {
						atomic.AddInt32(&queue.waitNum, -1)
						queue.lock.Unlock()
						queue.flag <- false
					} else {
						queue.lock.Unlock()
					}
				}
			}()
			flag := <-queue.flag
			tm.end()
			freeTimer(tm)
			if !flag {
				return nil, 0
			}
		} else {
			queue.lock.Unlock()
			v, num, has := queue.poll()
			if has {
				return v, num
			}
		}
	}
}

func NewVQueue() *VQueue {
	var ca int32 = 128
	el := newEle(ca)
	return &VQueue{write: el, read: el, flag: make(chan bool), num: 0, cap: ca, lock: new(sync.Mutex)}
}
