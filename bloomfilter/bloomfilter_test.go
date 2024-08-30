package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"sync"
	"testing"
)

func TestFPP(t *testing.T) {
	bf := New(100000, 0.01)
	t.Log(bf.k)
	t.Log(bf.m)

	// all added must contains
	var beg = 0
	for i := 0; i < 100; i++ {
		cur := rand.IntN(10000)
		var all = 0
		for j := beg; j < cur+beg; j++ {
			data := []byte(fmt.Sprintf("ele-%d", j))
			bf.Add(data)
			if bf.Contains(data) {
				all++
			}
		}
		assert.Equal(t, all, cur)
		beg += cur
	}

	// check fpp
	// 0.1    0.10354557965583433    0.10381663493587558
	// 0.01   0.010432378591592455   0.010646673785485943
	// 0.001  0.0010745359597678775  0011024807337292007
	// 0.0001 0.00013171636865849542 0.0001459728966773163
	bf = New(100000, 0.01)
	pool := make([]int, 0)
	for i := 0; i < 100000; i++ {
		pool = append(pool, i)
		bf.Add([]byte(fmt.Sprintf("ele-%d", i)))
	}
	var totalMiss = 0.0
	for i := 0; i < 100; i++ {
		var miss = 0
		cur := rand.IntN(100000)
		for j := 100000; j < cur+100000; j++ {
			data := []byte(fmt.Sprintf("ele-%d", j))
			if bf.Contains(data) {
				miss++
			}
		}
		totalMiss += float64(miss) / float64(cur)
	}
	fmt.Println(totalMiss / 100)
}

func TestMultiThreadSafe(t *testing.T) {
	bf := New(10000000, 0.01)
	t.Log(bf.k)
	t.Log(bf.m)
	wait := sync.WaitGroup{}
	wait.Add(1000)
	for j := 0; j < 1000; j++ {
		go func() {
			defer wait.Done()
			var beg = 0
			for i := 0; i < 1000; i++ {
				cur := rand.IntN(100)
				var all = 0
				for k := beg; k < cur+beg; k++ {
					data := []byte(fmt.Sprintf("ele-%d", k))
					bf.Add(data)
					if bf.Contains(data) {
						all++
					}
				}
				assert.Equal(t, cur, all)
				beg += cur
			}
		}()
	}
	wait.Wait()
}
func TestMarshal(t *testing.T) {
	bf := New(100000, 0.01)
	fmt.Println(bf.m)
	for i := 0; i < 100000; i++ {
		data := []byte(fmt.Sprintf("ele-%d", i))
		bf.Add(data)
	}
	bys, err := bf.Marshal()
	assert.Nil(t, err)
	var nbf BloomFilter
	err = nbf.Unmarshal(bys)
	for i := 0; i < 100000; i++ {
		data := []byte(fmt.Sprintf("ele-%d", i))
		assert.True(t, nbf.Contains(data))
	}

}
