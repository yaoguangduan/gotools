package mpsc

import (
	"fmt"
	"gotools/unbounded"
	"sync"
	"testing"
	"time"
)

// Benchmark for MPSC queue
func BenchmarkMPSC(b *testing.B) {
	ub := New[int]()

	go func() {
		for {
			ub.Poll()
		}
	}()
	w := sync.WaitGroup{}
	cur := time.Now().UnixMilli()
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for i := 0; i < 10000; i++ {
				ub.Add(i)
			}
		}()
	}
	w.Wait()
	fmt.Println(time.Now().UnixMilli() - cur)
}

func BenchmarkUnboundQueue(b *testing.B) {
	ub := unbounded.New[*int]()
	go func() {
		for {
			ub.Poll()
		}
	}()
	w := sync.WaitGroup{}
	cur := time.Now().UnixMilli()
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for i := 0; i < 10000; i++ {
				ub.Offer(&i)
			}
		}()
	}
	w.Wait()
	fmt.Println(time.Now().UnixMilli() - cur)
}
