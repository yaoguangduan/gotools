package mpsc

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
)

const (
	mask uint32 = 1023 * 2
)

type buffer[T any] struct {
	data []atomic.Pointer[T]
	next atomic.Pointer[buffer[T]]
}

func newBuffer[T any]() *buffer[T] {
	newChunk := make([]atomic.Pointer[T], 1024)
	buf := buffer[T]{data: newChunk}
	buf.next.Store(nil)
	return &buf
}

type MPSC[T any] struct {
	pIdx   atomic.Uint32
	_      cpu.CacheLinePad
	pLimit atomic.Uint32
	_      cpu.CacheLinePad
	pChunk *buffer[T]
	_      cpu.CacheLinePad
	cIdx   atomic.Uint32
	_      cpu.CacheLinePad
	cChunk *buffer[T]
	jump   *T
}

func New[T any]() *MPSC[T] {
	buf := newBuffer[T]()
	ub := &MPSC[T]{}
	ub.jump = new(T)
	ub.pChunk = buf
	ub.cChunk = buf
	ub.pLimit.Store(7)
	return ub
}

func (ub *MPSC[T]) Add(t T) {
	var idx uint32
	var chunk *buffer[T]
	for {
		limit := ub.pLimit.Load()
		idx = ub.pIdx.Load()
		if idx&1 == 1 {
			continue
		}
		chunk = ub.pChunk
		if limit <= idx {
			ope := ub.offerSlowPath(idx, limit)
			switch ope {
			case continueToPCas:
				break
			case retry:
				continue
			case resize:
				ub.resize(chunk, idx, t)
				return
			}
		}
		if ub.pIdx.CompareAndSwap(idx, idx+2) {
			break
		}
	}
	offset := (idx & mask) >> 1
	chunk.data[offset].Store(&t)
}

func (ub *MPSC[T]) Poll() *T {
	chunk := ub.cChunk
	idx := ub.cIdx.Load()
	offset := (idx & mask) >> 1
	var ptr = chunk.data[offset].Load()
	if ptr == nil {
		if idx != ub.pIdx.Load() {
			for {
				ptr = chunk.data[offset].Load()
				if ptr != nil {
					break
				}
			}
		} else {
			return nil
		}
	}
	if ptr == ub.jump {
		nc := ub.nextChunk(chunk)
		if nc.data == nil {
			panic("slice nl")
		}
		ptr = nc.data[offset].Load()
		if ptr == nil {
			panic("next buffer must have element")
			return nil
		} else {
			nc.data[offset].Store(nil)
			chunk.next.Store(nil)
			ub.cIdx.Store(idx + 2)
			return ptr
		}
	}
	chunk.data[offset].Store(nil)
	ub.cIdx.Store(idx + 2)
	return ptr
}

func (ub *MPSC[T]) Len() int {
	after := ub.cIdx.Load()
	var size int
	for {
		before := after
		pi := ub.pIdx.Load()
		after = ub.cIdx.Load()
		if before == after {
			size = int((pi - after) >> 1)
			break
		}
	}
	return size
}

func (ub *MPSC[T]) nextChunk(old *buffer[T]) *buffer[T] {
	ptrLast := old.next.Load()
	ub.cChunk = ptrLast
	return ptrLast
}

func (ub *MPSC[T]) offerSlowPath(idx uint32, limit uint32) op {
	cIdx := ub.cIdx.Load()
	capacity := mask
	if cIdx+capacity > idx {
		if ub.pLimit.CompareAndSwap(limit, cIdx+capacity) {
			return continueToPCas
		} else {
			return retry
		}
	} else if ub.pIdx.CompareAndSwap(idx, idx+1) {
		return resize
	} else {
		return retry
	}
}

func (ub *MPSC[T]) resize(old *buffer[T], idx uint32, t T) {
	offset := (idx & mask) >> 1
	newChunk := newBuffer[T]()
	if newChunk == nil {
		panic("cannot nil")
	}
	old.next.Store(newChunk)
	ub.pChunk = newChunk

	newChunk.data[offset].Store(&t)

	ub.pLimit.Store(idx + mask)

	ub.pIdx.Store(idx + 2)
	old.data[offset].Store(ub.jump)
}

type op uint8

const (
	continueToPCas op = iota
	retry
	resize
)
