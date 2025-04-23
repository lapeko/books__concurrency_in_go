package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mux sync.Mutex
	second := time.Second

	politeRun := func() {
		defer wg.Done()
		var counter int
		for startTime := time.Now(); time.Since(startTime) < second; {
			mux.Lock()
			time.Sleep(3 * time.Nanosecond)
			mux.Unlock()
			counter++
		}
		fmt.Printf("Greed run count: %d\n", counter)
	}

	greedyRun := func() {
		defer wg.Done()
		var counter int
		for startTime := time.Now(); time.Since(startTime) < second; {
			for count := 0; count < 3; count++ {
				mux.Lock()
				time.Sleep(time.Nanosecond)
				mux.Unlock()
			}
			counter++
		}
		fmt.Printf("Greed run count: %d\n", counter)
	}

	wg.Add(2)
	go politeRun()
	go greedyRun()
	wg.Wait()
}
