package lockfreequeue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLen(t *testing.T) {
	q := New[int]()
	wait := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for j := 0; j < 100; j++ {
				q.Enqueue(j)
			}
		}()
	}
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for j := 0; j < 10; j++ {
				_, b := q.Dequeue()
				if !b {
					j--
				}
			}
		}()
	}
	wait.Wait()
	fmt.Println(q.len.Load())
	assert.Equal(t, uint64(9000), q.Len())
}
