package panic

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"testing"
)

func Panic()  {
	if err := recover(); err != nil {
		fmt.Println(string(debug.Stack()))
	}
}

func TestLog(t *testing.T) {
	defer Panic()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
