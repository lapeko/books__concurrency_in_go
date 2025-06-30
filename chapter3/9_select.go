package main

import "fmt"

func main() {
	ch1, ch2 := make(chan struct{}), make(chan struct{})
	close(ch1)
	close(ch2)
	count, iterations := 0, 1000
	for i := 0; i < iterations; i++ {
		select {
		case <-ch1:
			count++
		case <-ch2:
		}
		fmt.Println("count1", count)
		fmt.Println("count2", iterations-count)
	}
}
