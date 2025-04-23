package main

import (
	"fmt"
	"sync"
)

func main() {
	var num int
	var mux sync.Mutex

	go func() {
		mux.Lock()
		num++
		mux.Unlock()
	}()

	mux.Lock()
	if num == 0 {
		fmt.Printf("num should be 0. Actual value: %d\n", num)
	} else {
		fmt.Printf("num should be 1. Actual value: %d\n", num)
	}
	mux.Unlock()
}
