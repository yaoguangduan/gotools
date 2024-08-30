package i64adder

import (
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
)

func TestInt64adder(t *testing.T) {
	adder := New()
	var w sync.WaitGroup
	w.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			adder.Incr()
			w.Done()
		}()
	}
	w.Wait()
	if adder.Sum() != 1000 {
		t.Errorf("Expected 1000, but got %d", adder.Sum())
	}
}

func TestAtomicInt64(t *testing.T) {
	var counter atomic.Int64
	var w sync.WaitGroup
	w.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			counter.Add(1)
			w.Done()
		}()
	}
	w.Wait()
	if counter.Load() != 1000 {
		t.Errorf("Expected 1000, but got %d", counter.Load())
	}
}

func TestRandom(t *testing.T) {
	adder := New()
	var sum atomic.Int64
	wait := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for j := 0; j < 100; j++ {
				val := rand.IntN(10000)
				sum.Add(int64(val))
				adder.Add(int64(val))
			}
		}()
	}
	wait.Wait()
	assert.Equal(t, sum.Load(), adder.Sum())
}
