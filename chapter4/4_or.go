package main

import (
	"fmt"
	"time"
)

func or(chans ...<-chan struct{}) <-chan struct{} {
	done := make(chan struct{})

	switch len(chans) {
	case 0:
		close(done)
	case 1:
		return chans[0]
	}

	go func() {
		defer close(done)
		select {
		case <-chans[0]:
		case <-chans[1]:
		case <-or(append(chans[2:], done)...):
		}
	}()

	return done
}

func main() {
	done1, done2, done3, done4 := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() { time.Sleep(time.Second * 2); done1 <- struct{}{} }()
	go func() { time.Sleep(time.Second); close(done2) }()
	go func() { time.Sleep(time.Second * 3); done3 <- struct{}{} }()
	go func() { time.Sleep(time.Second / 2); done3 <- struct{}{} }()
	start := time.Now()
	<-or(done1, done2, done3, done4)
	fmt.Printf("Done after %s\n", time.Since(start))
}
