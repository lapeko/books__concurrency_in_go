package main

import (
	"fmt"
	"sync"
	"time"
)

type Button struct {
	Clicked *sync.Cond
}

func subscribe(b *Button, f func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer b.Clicked.L.Unlock()
		b.Clicked.L.Lock()
		wg.Done()
		b.Clicked.Wait()
		f()
	}()
	wg.Wait()
}

func main() {
	clicked := sync.NewCond(&sync.Mutex{})
	subscribe(&Button{Clicked: clicked}, func() {
		fmt.Println("Full size window clicked")
	})
	subscribe(&Button{Clicked: clicked}, func() {
		fmt.Println("Cancel button clicked")
	})
	subscribe(&Button{Clicked: clicked}, func() {
		fmt.Println("Ok button clicked")
	})
	clicked.Broadcast()
	time.Sleep(time.Millisecond)
}
