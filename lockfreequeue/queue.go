package lockfreequeue

import (
	"sync/atomic"
	"unsafe"
)

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

type LockFreeQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	len  atomic.Uint64
	zero T
}

func New[T any]() *LockFreeQueue[T] {
	n := unsafe.Pointer(&node[T]{})
	q := &LockFreeQueue[T]{head: n, tail: n}
	q.zero = *new(T)
	return q
}
func (q *LockFreeQueue[T]) Len() uint64 {
	return q.len.Load()
}

func (q *LockFreeQueue[T]) Enqueue(value T) {
	n := unsafe.Pointer(&node[T]{value: value})
	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*node[T])(tail).next)
		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&(*node[T])(tail).next, next, n) {
					atomic.CompareAndSwapPointer(&q.tail, tail, n)
					q.len.Add(1)
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*node[T])(head).next)
		if head == atomic.LoadPointer(&q.head) {
			if head == tail {
				if next == nil {
					return q.zero, false
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				value := (*node[T])(next).value
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					q.len.Add(^uint64(0))
					return value, true
				}
			}
		}
	}
}
