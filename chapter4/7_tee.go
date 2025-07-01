package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func tee(in <-chan int) (<-chan int, <-chan int) {
	out1, out2 := make(chan int), make(chan int)

	go func() {
		defer func() {
			close(out1)
			close(out2)
		}()

		for num := range in {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() { out1 <- num; wg.Done() }()
			go func() { out2 <- num; wg.Done() }()
			wg.Wait()
		}
	}()

	return out1, out2
}

func generateNInts(size int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := 1; num <= size; num++ {
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(500)))
			out <- num
		}
	}()

	return out
}

func main() {
	size := 10
	out1, out2 := tee(generateNInts(size))
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range out1 {
			fmt.Println("out1:", num)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range out2 {
			fmt.Println("out2:", num)
		}
	}()

	wg.Wait()
}
