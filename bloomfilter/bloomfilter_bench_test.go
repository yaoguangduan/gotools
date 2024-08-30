package algo

import (
	"math/rand"
	"testing"
	"time"
)

/**
Bloom VS bits-and-blooms/bloom/v3

BenchmarkBloomFilter_Add-8        	16167429	        70.44 ns/op
BenchmarkBloomFilter_Contains-8   	 4327959	       269.0 ns/op

BenchmarkBitsAndBloomsBloom_Add-8        	11760151	        93.56 ns/op
BenchmarkBitsAndBloomsBloom_Contains-8   	 4069503	       301.3 ns/op
*/

// 基准测试 Add 接口
func BenchmarkBloomFilter_Add(b *testing.B) {
	bf := NewWithInsertion(1000000)
	rand.Seed(time.Now().UnixNano())
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = make([]byte, 16)
		rand.Read(data[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add(data[i])
	}
}

// 基准测试 Contains 接口
func BenchmarkBloomFilter_Contains(b *testing.B) {
	bf := NewWithInsertion(1000000)
	rand.Seed(time.Now().UnixNano())
	data := make([][]byte, 1000000)
	for i := 0; i < 1000000; i++ {
		data[i] = make([]byte, 16)
		rand.Read(data[i])
		bf.Add(data[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Contains(data[rand.Intn(1000000)])
	}
}
