package log

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func Test_timer(t *testing.T) {

	var iii int32 = 0
	for i := 0; i < 5_000_00; i++ {
		tm:=time.NewTimer(time.Second*10)
		go func() {
			<-tm.C
			atomic.AddInt32(&iii,1)
			Info(iii)
		}()
	}




	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig

}