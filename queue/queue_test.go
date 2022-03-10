package queue

import (
	"log"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	lite:=NewLiteQueue(3)

	go func() {
		for{
			time.Sleep(time.Second)
			v,num:=lite.poll()
			log.Println(v,num)

		}
	}()

	lite.Offer("45")
	lite.Offer("12364456")
	t.Log("~~~~~~")
	lite.Offer("123464434456")
	t.Log("~~~~~~")
	lite.Offer("123465656")
	t.Log("~~~~~~")
	lite.Offer("123433356")
	t.Log("~~~~~~")
	time.Sleep(time.Hour)

}
