package queue

import (
	"context"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	lite:=NewQueue()


	go func() {
		for{
			time.Sleep(time.Second*10)
			lite.Offer("123")
		}
	}()


	for{
		ctx,  CancelFunc:= context.WithTimeout(context.Background(), time.Second*2)
		v,num,close := lite.Dequeue(ctx)
		if !close{
			CancelFunc()
		}
		t.Log(v,num,"close:",close)
	}


}
