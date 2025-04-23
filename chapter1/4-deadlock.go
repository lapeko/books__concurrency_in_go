package main

import (
	"sync"
	"time"
)

type value struct {
	mux sync.Mutex
}

func main() {
	var a, b value
	var wg = sync.WaitGroup{}

	fn := func(a, b *value) {
		defer wg.Done()
		defer a.mux.Unlock()
		a.mux.Lock()

		time.Sleep(time.Second)

		defer b.mux.Unlock()
		b.mux.Lock()
	}

	wg.Add(2)
	go fn(&a, &b)
	go fn(&b, &a)
	wg.Wait()
}
