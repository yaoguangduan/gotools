package unbounded

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"sync"
	"testing"
)

func TestMPMC(t *testing.T) {
	ub := New[int]()
	defer ub.Close()
	size := 1000
	wait := sync.WaitGroup{}
	for i := 0; i < size; i += 100 {
		wait.Add(1)
		go func(j int) {
			defer wait.Done()
			for k := j; k < j+100; k++ {
				ub.Offer(k)
			}
		}(i)
	}
	res := make([]int, 0)
	l := sync.Mutex{}
	for i := 0; i < size; i += 100 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for {
				l.Lock()
				if len(res) >= size {
					l.Unlock()
					return
				}
				val := ub.Poll()
				res = append(res, val)
				l.Unlock()
			}
		}()
	}
	wait.Wait()
	var allContain = true
	for i := range size {
		if !slices.Contains(res, i) {
			allContain = false
			fmt.Println(i)
			break
		}
	}
	if !allContain {
		fmt.Println(res)
		t.Errorf("invalid consume seq:%+v", res)
	}
}

func TestMPSC(t *testing.T) {
	ub := New[int]()
	defer ub.Close()
	size := 1000
	wait := sync.WaitGroup{}
	for i := 0; i < size; i += 100 {
		wait.Add(1)
		go func(j int) {
			defer wait.Done()
			for k := j; k < j+100; k++ {
				ub.Offer(k)
			}
		}(i)
	}
	res := make([]int, 0)
	wait.Add(1)
	go func() {
		defer wait.Done()
		for len(res) != size {
			res = append(res, ub.Poll())
		}
	}()
	wait.Wait()
	var allContain = true
	for i := range size {
		if !slices.Contains(res, i) {
			allContain = false
			fmt.Println(i)
			break
		}
	}
	if !allContain {
		fmt.Println(res)
		t.Errorf("invalid consume seq:%+v", res)
	}
}

func TestSPSC(t *testing.T) {
	ub := New[int]()
	defer ub.Close()
	size := 1000
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		defer wait.Done()
		for k := 0; k < size; k++ {
			ub.Offer(k)
		}
	}()
	res := make([]int, 0)
	wait.Add(1)
	go func() {
		defer wait.Done()
		for len(res) != size {
			res = append(res, ub.Poll())
		}
	}()
	wait.Wait()
	var allContain = true
	for i := range size {
		if !slices.Contains(res, i) {
			allContain = false
			fmt.Println(i)
			break
		}
	}
	if !allContain {
		fmt.Println(res)
		t.Errorf("invalid consume seq:%+v", res)
	}
}
func TestManyPutGet(t *testing.T) {
	for range 1000 {
		TestMPSC(t)
		TestSPSC(t)
		TestMPMC(t)
	}
}

func TestClose(t *testing.T) {
	ub := New[int]()
	ub.Offer(1)
	ub.Offer(2)
	ub.Offer(3)
	fmt.Println(ub)
	fmt.Println(ub.Len())
	ub.Close()
	fmt.Println(ub.Poll())
	fmt.Println(ub.Poll())
	fmt.Println(ub.Poll())
	assert.Equal(t, ub.Len(), 0)
	assert.Equal(t, ub.Poll(), 0)
	ub = New[int]()
	ub.Offer(4)
	ub.Offer(5)
	ub.Close()
	fmt.Println("iter")
	assert.Equal(t, []int{4, 5}, slices.Collect(ub.Iter()))
}
