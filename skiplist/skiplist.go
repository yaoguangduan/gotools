package algo

import (
	"cmp"
	"iter"
	"math"
	"math/rand/v2"
	"time"
)

type CompareFunc[T any] func(a, b T) int

const maxLevel = 32

type skipListNode[K any, V any] struct {
	key   K
	value V
	next  []*skipListNode[K, V]
}

type SkipList[K any, V any] struct {
	head        *skipListNode[K, V]
	cache       []*skipListNode[K, V]
	level       int
	cmpFn       CompareFunc[K]
	rand        *rand.Rand
	len         int
	nodeInitFns []func(K, V) *skipListNode[K, V]
}

func NewSkipList[K cmp.Ordered, V any]() *SkipList[K, V] {
	return NewSkipListWithCmp[K, V](func(a, b K) int {
		return cmp.Compare(a, b)
	})
}
func NewSkipListWithCmp[K any, V any](cmpFn CompareFunc[K]) *SkipList[K, V] {
	slm := &SkipList[K, V]{rand: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 1)), cmpFn: cmpFn,
		head:        &skipListNode[K, V]{next: make([]*skipListNode[K, V], maxLevel)},
		cache:       make([]*skipListNode[K, V], maxLevel),
		level:       1,
		nodeInitFns: make([]func(K, V) *skipListNode[K, V], maxLevel+1),
	}
	slm.initNodeInitFn()
	return slm
}

func (s *SkipList[K, V]) Clear() {
	for i := range s.head.next {
		s.head.next[i] = nil
	}
	s.level = 1
	s.len = 0
}

func (s *SkipList[K, V]) Len() int {
	return s.len
}
func (s *SkipList[K, V]) Has(key K) bool {
	find, _ := s.findByKeyOrPrev(key)
	return find != nil
}
func (s *SkipList[K, V]) Iter() iter.Seq[K] {
	return func(yield func(K) bool) {
		for tmp := s.head.next[0]; tmp != nil; tmp = tmp.next[0] {
			if !yield(tmp.key) {
				break
			}
		}
	}
}
func (s *SkipList[K, V]) Iter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for tmp := s.head.next[0]; tmp != nil; tmp = tmp.next[0] {
			if !yield(tmp.key, tmp.value) {
				break
			}
		}
	}
}
func (s *SkipList[K, V]) IterRange(beg, end K) iter.Seq2[K, V] {
	bf, prev := s.findByKeyAndPrev(beg)
	if bf == nil {
		bf = prev[0]
	}
	return func(yield func(K, V) bool) {
		for tmp := bf; s.cmpFn(tmp.key, end) < 0; tmp = tmp.next[0] {
			if !yield(tmp.key, tmp.value) {
				break
			}
		}
	}
}
func (s *SkipList[K, V]) IterBE(beg K) iter.Seq2[K, V] {
	bf, prev := s.findByKeyAndPrev(beg)
	if bf == nil {
		bf = prev[0].next[0]
	}
	return func(yield func(K, V) bool) {
		for tmp := bf; tmp != nil; tmp = tmp.next[0] {
			if !yield(tmp.key, tmp.value) {
				break
			}
		}
	}
}

func (s *SkipList[K, V]) IterLE(end K) iter.Seq2[K, V] {
	bf, prev := s.findByKeyAndPrev(end)
	if bf == nil {
		bf = prev[0]
	}
	return func(yield func(K, V) bool) {
		for tmp := s.head.next[0]; tmp != bf.next[0]; tmp = tmp.next[0] {
			if !yield(tmp.key, tmp.value) {
				break
			}
		}
	}
}

func (s *SkipList[K, V]) Put(key K, value V) {
	find, updates := s.findByKeyOrPrev(key)
	if find != nil {
		find.value = value
		return
	}
	level := s.genLevel()
	node := s.nodeInitFns[level](key, value)
	minLevel := min(level, s.level)
	for i := 0; i < minLevel; i++ {
		node.next[i] = updates[i].next[i]
		updates[i].next[i] = node
	}
	if level > s.level {
		for i := s.level; i < level; i++ {
			s.head.next[i] = node
		}
		s.level = level
	}
	s.len++
}
func (s *SkipList[K, V]) Get(key K) V {
	f, _ := s.findByKeyOrPrev(key)
	if f != nil {
		return f.value
	}
	return *new(V)
}

func (s *SkipList[K, V]) TryGet(key K) (V, bool) {
	f, _ := s.findByKeyOrPrev(key)
	if f != nil {
		return f.value, true
	}
	return *new(V), false
}

func (s *SkipList[K, V]) GetOr(key K, val V) V {
	f, _ := s.findByKeyOrPrev(key)
	if f != nil {
		return f.value
	}
	return val
}

func (s *SkipList[K, V]) Delete(key K) V {
	f, updates := s.findByKeyAndPrev(key)
	if f == nil {
		return *new(V)
	}
	for i, v := range f.next {
		updates[i].next[i] = v
	}

	for s.level > 1 && s.head.next[s.level-1] == nil {
		s.level--
	}
	s.len--
	return f.value
}

func (s *SkipList[K, V]) findByKeyOrPrev(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	updates := s.cache[0:s.level]
	prev := s.head
	for i := s.level - 1; i >= 0; i-- {
		for next := prev.next[i]; next != nil; next = next.next[i] {
			cr := s.cmpFn(next.key, key)
			if cr == 0 {
				return next, nil
			}
			if cr > 0 {
				break
			}
			prev = next
		}
		updates[i] = prev
	}
	return nil, updates
}
func (s *SkipList[K, V]) findByKeyAndPrev(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	updates := s.cache[0:s.level]
	prev := s.head
	var find *skipListNode[K, V]
	for i := s.level - 1; i >= 0; i-- {
		for next := prev; next != nil; next = next.next[i] {
			cr := s.cmpFn(next.key, key)
			if cr == 0 {
				find = next
				break
			}
			if cr > 0 {
				break
			}
			updates[i] = next
		}
	}
	return find, updates
}

const half = math.MaxUint64 / 2

func (s *SkipList[K, V]) genLevel() int {
	var level = 1
	for level < 32 && s.rand.Uint64() < half {
		level++
	}
	return level
}

func (s *SkipList[K, V]) initNodeInitFn() {
	s.nodeInitFns[1] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [1]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[2] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [2]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[3] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [3]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[4] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [4]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[5] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [5]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[6] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [6]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[7] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [7]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[8] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [8]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[9] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [9]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[10] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [10]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[11] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [11]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[12] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [12]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[13] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [13]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[14] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [14]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[15] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [15]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[16] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [16]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[17] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [17]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[18] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [18]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[19] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [19]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[20] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [20]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[21] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [21]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[22] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [22]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[23] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [23]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[24] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [24]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[25] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [25]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[26] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [26]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[27] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [27]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[28] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [28]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[29] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [29]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[30] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [30]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[31] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [31]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
	s.nodeInitFns[32] = func(k K, v V) *skipListNode[K, V] {
		n := struct {
			head  skipListNode[K, V]
			tails [32]*skipListNode[K, V]
		}{head: skipListNode[K, V]{k, v, nil}}
		n.head.next = n.tails[:]
		return &n.head
	}
}
