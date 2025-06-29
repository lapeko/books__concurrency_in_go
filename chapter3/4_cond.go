package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	size := 10
	wg.Add(size)
	list := make([]struct{}, 0, size)
	cond := sync.NewCond(&sync.Mutex{})

	pop := func() {
		defer func() {
			cond.Signal()
			cond.L.Unlock()
			wg.Done()
		}()

		cond.L.Lock()
		fmt.Println("Pop from list")
		list = list[1:]
	}

	for i := 0; i < size; i++ {
		cond.L.Lock()
		for len(list) >= 2 {
			cond.Wait()
		}
		fmt.Println("Pushing into list")
		list = append(list, struct{}{})
		cond.L.Unlock()

		go pop()
	}
	wg.Wait()
}
