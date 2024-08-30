package mpsc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"sync"
	"testing"
)

func TestOp(t *testing.T) {
	ub := New[int]()
	ub.Poll()
	ub.Add(1)
	ub.Add(2)
	ub.Add(3)
	ub.Add(4)
	ub.Add(5)
	ub.Add(6)
	ub.Add(7)
	ub.Add(8)
	ub.Add(9)
	ub.Add(10)
	ub.Add(11)
	ub.Add(12)
	ub.Add(13)
	ub.Add(14)
	ub.Add(15)
	ub.Add(16)
	ub.Add(17)
	ub.Add(18)
	ub.Add(19)
	fmt.Println(ub)
	for i := 1; i <= 19; i++ {
		assert.Equal(t, i, *ub.Poll())
	}
	fmt.Println(ub.Poll())
}
func TestAdd(t *testing.T) {
	ub := New[int]()
	w := sync.WaitGroup{}
	for i := 1; i < 100; i += 10 {
		w.Add(1)
		go func(j int) {
			defer w.Done()
			for k := j; k < j+10; k++ {
				ub.Add(k)
			}
		}(i)
	}
	w.Wait()
	res := make([]int, 0)

	chnk := ub.cChunk
	for {
		if chnk == nil {
			break
		}
		slice := chnk.data
		for i, _ := range slice {
			if slice[i].Load() == nil {
				continue
			} else {
				if *slice[i].Load() != 0 {
					vv := *slice[i].Load()
					if vv != 0 {
						res = append(res, vv)
					}
				}
			}
		}
		chnk = chnk.next.Load()
	}
	slices.Sort(res)
	var containAll bool = true
	for i := range 101 {
		if i == 0 {
			continue
		}
		if !slices.Contains(res, i) {
			containAll = false
		}
	}
	if len(res) != 100 || !containAll {
		fmt.Println(len(res))
		fmt.Println(res)
	}
}
func TestAddWithPoll(t *testing.T) {
	ub := New[int]()
	w := sync.WaitGroup{}
	w1 := sync.WaitGroup{}
	w1.Add(1)
	go func() {
		defer w1.Done()
		res := make([]int, 0)
		var val *int
		for {
			val = ub.Poll()
			if val != nil && *val == 999 {
				break
			}
			if val != nil {
				res = append(res, *val)
			}
		}

		slices.Sort(res)
		var containAll = true
		for i := range 101 {
			if i == 0 {
				continue
			}
			if !slices.Contains(res, i) {
				containAll = false
			}
		}
		if len(res) != 100 || !containAll {
			fmt.Println(len(res))
			fmt.Println(res)
		}
	}()
	for i := 1; i < 100; i += 10 {
		w.Add(1)
		go func(j int) {
			defer w.Done()
			for k := j; k < j+10; k++ {
				ub.Add(k)
			}
		}(i)
	}
	w.Wait()
	ub.Add(999)
	w1.Wait()
}
func TestManyAdd(t *testing.T) {

	for i := 0; i < 100000; i++ {
		TestAdd(t)
	}
}

func TestManyAddWithPoll(t *testing.T) {

	for i := 0; i < 100000; i++ {
		TestAddWithPoll(t)
	}
}
func TestSize(t *testing.T) {
	m := New[int]()
	m.Add(1)
	m.Add(2)
	m.Add(3)
	assert.Equal(t, 3, m.Len())
}
