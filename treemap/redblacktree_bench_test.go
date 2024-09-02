package algo

import (
	"math/rand"
	"testing"
	"time"
)

/**

BenchmarkRedBlackTreeRandOp-8   	 1797362	       726.3 ns/op
BenchmarkRedBlackTreeRandOp-8   	 1980540	       782.9 ns/op
BenchmarkRedBlackTreeRandOp-8   	 2040639	       738.9 ns/op

optimize:
BenchmarkRedBlackTreeRandOp-8   	 1930863	       696.1 ns/op
BenchmarkRedBlackTreeRandOp-8   	 2043625	       686.2 ns/op
BenchmarkRedBlackTreeRandOp-8   	 1863933	       691.3 ns/op
*/

func BenchmarkRedBlackTreeRandOp(b *testing.B) {
	tree := NewWithCmpFunc[int, int](func(a, b int) int {
		return a - b
	})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		op := rand.Intn(2)
		if op == 0 {
			tree.Get(rand.Intn(b.N))
		} else {
			tree.Put(rand.Intn(b.N), i)
		}
	}
}

func BenchmarkRedBlackTreeInsert(b *testing.B) {
	tree := NewWithCmpFunc[int, int](func(a, b int) int {
		return a - b
	})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		tree.Put(rand.Intn(b.N), i)
	}
}

// 基准测试红黑树查找
func BenchmarkRedBlackTreeGet(b *testing.B) {
	tree := NewWithCmpFunc[int, int](func(a, b int) int {
		return a - b
	})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		tree.Put(rand.Intn(b.N), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Get(rand.Intn(b.N))
	}
}

// 基准测试内置 map 插入
func BenchmarkMapInsert(b *testing.B) {
	m := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		m[rand.Intn(b.N)] = i
	}
}

// 基准测试内置 map 查找
func BenchmarkMapGet(b *testing.B) {
	m := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		m[rand.Intn(b.N)] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[rand.Intn(b.N)]
	}
}
