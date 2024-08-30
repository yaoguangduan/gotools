package i64adder

import (
	"sync"
	"sync/atomic"
	"testing"
)

/*
*
cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
BenchmarkLockMutex
BenchmarkLockMutex-8     	   10000	    106789 ns/op
BenchmarkU64Adder
BenchmarkU64Adder-8      	  588508	      2389 ns/op
BenchmarkAtomicInt64
BenchmarkAtomicInt64-8   	   62172	     19311 ns/op
*/
func BenchmarkLockMutex(b *testing.B) {
	var counter = 0
	lock := sync.Mutex{}
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				lock.Lock()
				counter++
				lock.Unlock()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkU64Adder(b *testing.B) {
	adder := New()
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				adder.Add(1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkAtomicInt64(b *testing.B) {
	var counter atomic.Int64
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {

				counter.Add(1)
			}
		}()
	}
	wg.Wait()
}
