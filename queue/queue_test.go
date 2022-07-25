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
			time.Sleep(time.Second*2)
			lite.Offer("123")
		}
	}()


	for{
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		v,num,close := lite.Dequeue(ctx)
		t.Log(v,num,"close:",close)
	}


}
