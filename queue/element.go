package queue

import (
	"sync"
)

type element struct {
	next  *element
	value interface{}
}

func newElement(value interface{}) *element {
	return &element{value: value}
}
var poolElement = &sync.Pool{
	New: func() interface{} {
		return new(element)
	},
}

func getElement(value interface{}) *element {
	ele:= poolElement.Get().(*element)
	ele.value = value
	return ele
}
func freeElement(ele *element) {
	ele.next = nil
	poolElement.Put(ele)
}

