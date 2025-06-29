package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mut sync.Mutex
	var num int
	size := 5

	inc := func() {
		defer mut.Unlock()
		mut.Lock()
		num++
		fmt.Printf("inc(); num: %d\n", num)
	}

	dec := func() {
		defer mut.Unlock()
		mut.Lock()
		num--
		fmt.Printf("dec(); num: %d\n", num)
	}

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			inc()
		}()
	}

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			dec()
		}()
	}

	wg.Wait()
}
