package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mux   sync.Mutex
	count int
}

func (c *Counter) inc() {
	defer c.mux.Unlock()
	c.mux.Lock()
	c.count++
}

func main() {
	size := 100
	var wg sync.WaitGroup
	var once sync.Once
	var counter Counter

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			once.Do(counter.inc)
		}()
	}
	wg.Wait()

	fmt.Println(counter.count)
}
