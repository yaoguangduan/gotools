package i64adder

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
)

const adderChunkSize = 16

type cell struct {
	_ cpu.CacheLinePad
	n int64
	_ cpu.CacheLinePad
}

type Adder struct {
	cells []cell
}

func New() *Adder {
	c := &Adder{cells: make([]cell, adderChunkSize)}
	for i := range c.cells {
		c.cells[i] = cell{}
	}
	return c
}
func (addr *Adder) Add(x int64) {
	idx := hash() & (adderChunkSize - 1)
	atomic.AddInt64(&addr.cells[idx].n, x)
}
func (addr *Adder) Incr() {
	addr.Add(1)
}
func (addr *Adder) Decr() {
	addr.Add(-1)
}
func (addr *Adder) Sum() int64 {
	sum := int64(0)
	for _, n := range addr.cells {
		sum += n.n
	}
	return sum
}
