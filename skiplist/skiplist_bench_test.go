package algo

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

const (
	benchInitSize  = 1000000
	benchBatchSize = 10
)

func newSkipListN(n int) *SkipList[int, int] {
	sl := NewSkipList[int, int]()
	sl.rand = rand.New(rand.NewPCG(0, 1))
	for i := 0; i < n; i++ {
		sl.Put(i, i)
	}
	return sl
}

func newMapN(n int) map[int]int {
	m := map[int]int{}
	for i := 0; i < n; i++ {
		m[i] = i
	}
	return m
}

func BenchmarkSkipList_Iterate(b *testing.B) {
	sl := newSkipListN(100)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for k, v := range sl.Iter2() {
			_, _ = k, v
		}
	}
}

func BenchmarkSkipList_Put(b *testing.B) {
	start := benchInitSize
	sl := newSkipListN(start)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sl.Put(start+1, 1)
		start += benchBatchSize
	}
}

func BenchmarkMap_Put(b *testing.B) {
	start := benchInitSize
	m := newMapN(start)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchBatchSize; i++ {
			m[start+i] = i
		}
		start += benchBatchSize
	}
}

func BenchmarkSkipList_Put_Dup(b *testing.B) {
	sl := newSkipListN(benchInitSize)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchBatchSize; i++ {
			sl.Put(i, i)
		}
	}
}

func BenchmarkMap_Put_Dup(b *testing.B) {
	m := newMapN(benchInitSize)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchBatchSize; i++ {
			m[i] = i
		}
	}
}

func BenchmarkMap_Find(b *testing.B) {
	m := newMapN(benchInitSize)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchBatchSize; i++ {
			_, _ = m[i]
		}
	}
}

func BenchmarkSkipList_Find(b *testing.B) {
	sl := newSkipListN(benchInitSize)
	b.ResetTimer()
	b.Run("Find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				_ = sl.Get(n)
			}
		}
	})
	//b.Run("LowerBound", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		for n := 0; n < benchBatchSize; n++ {
	//			_ = sl.impl.lowerBound(n)
	//		}
	//	}
	//})
	//b.Run("FindEnd", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		for n := 0; n < benchBatchSize; n++ {
	//			_ = sl.Find(benchInitSize)
	//		}
	//	}
	//})
}

func BenchmarkSkipListString(b *testing.B) {
	sl := NewSkipList[string, int]()
	sl.rand = rand.New(rand.NewPCG(0, 1))
	var a []string
	for i := 0; i < benchBatchSize; i++ {
		a = append(a, strconv.Itoa(benchInitSize+i))
	}
	end := strconv.Itoa(2 * benchInitSize)
	b.ResetTimer()
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				sl.Put(a[n], n)
			}
		}
	})
	b.Run("Find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				sl.Get(a[n])
			}
		}
	})
	b.Run("FindEnd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				sl.Get(end)
			}
		}
	})

	b.Run("RemoveEnd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				sl.Delete(end)
			}
		}
	})
	b.Run("Remove", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for n := 0; n < benchBatchSize; n++ {
				sl.Delete(a[n])
			}
		}
	})
}
