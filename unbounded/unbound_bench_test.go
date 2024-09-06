package unbounded

import (
	"fmt"
	"github.com/puzpuzpuz/xsync/v3"
	"sync"
	"testing"
	"time"
)

func BenchmarkIntChan(b *testing.B) {
	b.SetParallelism(5)
	ch := make(chan *int, 128)
	go func() {
		for {
			<-ch
		}
	}()
	w := sync.WaitGroup{}
	cur := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for i := 0; i < 1000000; i++ {
				ch <- &i
			}
		}()
	}
	w.Wait()
	fmt.Println(time.Now().UnixMilli() - cur)
}

func BenchmarkUnboundQueue(b *testing.B) {
	b.SetParallelism(5)
	ub := New[*int]()
	go func() {
		for {
			ub.Poll()
		}
	}()
	w := sync.WaitGroup{}
	cur := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for i := 0; i < 1000000; i++ {
				ub.Offer(&i)
			}
		}()
	}
	w.Wait()
	ub.Close()
	fmt.Println(time.Now().UnixMilli() - cur)
}

func BenchmarkXSyncQueue(b *testing.B) {
	b.SetParallelism(5)
	ub := xsync.NewMPMCQueueOf[int](128)
	go func() {
		for {
			ub.Dequeue()
		}
	}()
	w := sync.WaitGroup{}
	cur := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for i := 0; i < 1000000; i++ {
				ub.Enqueue(i)
			}
		}()
	}
	w.Wait()
	fmt.Println(time.Now().UnixMilli() - cur)
}
