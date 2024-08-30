package unbounded

import (
	"fmt"
	"iter"
	"strings"
	"sync"
	"sync/atomic"
)

const cacheSegmentSize = 8
const cacheIdxMask = 7

type cache[T any] struct {
	data    []T
	jumpIdx int
	next    *cache[T]
}

func (c *cache[T]) init() {

}

type Chan[T any] struct {
	in     chan T
	out    chan T
	pCache *cache[T]
	cCache *cache[T]
	pLimit int
	pIdx   int
	cIdx   int
	close  atomic.Bool
	pool   sync.Pool
}

func New[T any]() *Chan[T] {
	ub := &Chan[T]{in: make(chan T, 1), out: make(chan T, 1)}
	go ub.start()
	return ub
}
func (u *Chan[T]) Offer(item T) {
	u.in <- item
}
func (u *Chan[T]) Poll() T {
	return <-u.out
}
func (u *Chan[T]) TryPoll() (T, bool) {
	v, b := <-u.out
	return v, b
}
func (u *Chan[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, o := <-u.out
			if !o || !yield(v) {
				break
			}
		}
	}
}
func (u *Chan[T]) Len() int {
	return len(u.out) + (u.pIdx - u.cIdx)
}

func (u *Chan[T]) Close() {
	if !u.close.CompareAndSwap(false, true) {
		panic("close a closed unbounded")
	}
	close(u.in)
}
func (u *Chan[T]) start() {
	defer close(u.out)
	u.initCache()
loop:
	for {
		nv, ok := <-u.in
		if !ok {
			break loop
		}
		select {
		case u.out <- nv:
			continue
		default:
		}
		u.intoCache(nv)
		var v = u.seeCached()
		for v != nil {
			select {
			case newVal, okNow := <-u.in:
				if !okNow {
					break loop
				} else {
					u.intoCache(newVal)
				}
			case u.out <- *v:
				u.consumeSeed()
				v = u.seeCached()
			}
		}
	}
	u.drain()
}

func (u *Chan[T]) initCache() {
	u.pool.New = func() interface{} {
		c := new(cache[T])
		c.data = make([]T, cacheSegmentSize)
		c.jumpIdx = -1
		return c
	}
	c := u.pool.Get().(*cache[T])
	u.pLimit = cacheIdxMask
	u.pCache, u.cCache = c, c
}

func (u *Chan[T]) seeCached() *T {
	if u.cIdx == u.pIdx {
		return nil
	}
	offset := u.cIdx & cacheIdxMask
	if offset == u.cCache.jumpIdx {
		tmp := u.cCache
		u.cCache = u.cCache.next
		tmp.next = nil
		tmp.jumpIdx = -1
		u.pool.Put(tmp)
	}
	val := &u.cCache.data[offset]
	return val
}

func (u *Chan[T]) drain() {
	var val = u.seeCached()
	for val != nil {
		u.consumeSeed()
		u.out <- *val
		val = u.seeCached()
	}
}

func (u *Chan[T]) intoCache(nv T) {
	offset := u.pIdx & cacheIdxMask
	if u.pIdx >= u.pLimit {
		if u.cIdx+cacheIdxMask > u.pIdx {
			u.pLimit = u.cIdx + cacheIdxMask
		} else {
			u.pCache.jumpIdx = offset
			c := u.pool.Get().(*cache[T])
			u.pCache.next = c
			u.pCache = c
			u.pLimit += cacheIdxMask
		}
	}
	u.pCache.data[offset] = nv
	u.pIdx++
}

func (u *Chan[T]) String() string {
	sb := strings.Builder{}
	var tmp = u.cCache
	for tmp != nil {
		sb.WriteString(fmt.Sprintf("%+v", tmp.data))
		sb.WriteString("\r\n")
		tmp = tmp.next
	}
	return sb.String()
}

func (u *Chan[T]) consumeSeed() {
	u.cIdx++
}
