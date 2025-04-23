package main

import "fmt"

func main() {
	var num int

	go func() {
		num++
	}()

	if num == 0 {
		fmt.Printf("num should be 0. Actual value: %d\n", num)
	} else {
		fmt.Printf("num should be 1. Actual value: %d\n", num)
	}
}
