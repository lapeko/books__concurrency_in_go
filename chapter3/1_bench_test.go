package main

import (
	"sync"
	"testing"
)

func BenchmarkMyContextSwitch(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	startCh, senderReceiver := make(chan struct{}), make(chan struct{})

	runSender := func() {
		defer wg.Done()
		<-startCh
		for i := 0; i < b.N; i++ {
			senderReceiver <- struct{}{}
		}
	}

	runReceiver := func() {
		defer wg.Done()
		<-startCh
		for i := 0; i < b.N; i++ {
			<-senderReceiver
		}
	}

	go runSender()
	go runReceiver()

	b.ResetTimer()
	close(startCh)
	wg.Wait()
}
