package main

import (
	"sync"
	"testing"
)

func runSender(start <-chan struct{}, send chan<- struct{}, wg *sync.WaitGroup, size int) {
	defer wg.Done()
	<-start
	for i := 0; i < size; i++ {
		send <- struct{}{}
	}
}

func runReceiver(start <-chan struct{}, receive <-chan struct{}, wg *sync.WaitGroup, size int) {
	defer wg.Done()
	<-start
	for i := 0; i < size; i++ {
		<-receive
	}
}

func BenchmarkMyContextSwitch(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	start, senderReceiver := make(chan struct{}), make(chan struct{})

	go runSender(start, senderReceiver, &wg, b.N)
	go runReceiver(start, senderReceiver, &wg, b.N)

	b.ResetTimer()
	close(start)
	wg.Wait()
}
