package main

import (
	"fmt"
	"sync"
)

func main() {
	var num int
	var mux sync.Mutex
	var ch = make(chan interface{})

	go func() {
		mux.Lock()
		num++
		mux.Unlock()
		ch <- struct{}{}
	}()

	<-ch

	mux.Lock()
	if num == 0 {
		fmt.Printf("num should be 0. Actual value: %d\n", num)
	} else {
		fmt.Printf("num should be 1. Actual value: %d\n", num)
	}
	mux.Unlock()
}
