package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedBlackTreeOp(t *testing.T) {
	rbt := New[int, int]()
	rbt.Put(5, 5)
	rbt.Put(2, 22)
	rbt.Put(6, 6)
	rbt.Put(4, 44)
	rbt.Put(2, 2)
	rbt.Put(1, 11)
	rbt.Put(7, 7)
	rbt.Put(4, 4)
	rbt.Put(8, 8)
	rbt.Put(1, 1)
	rbt.Put(3, 3)
	tmp := make([]int, 0)
	for k, _ := range rbt.Iter() {
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, tmp)
	tmp = make([]int, 0)
	for k, _ := range rbt.IterRange(5, 10) {
		fmt.Println(k)
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{5, 6, 7, 8}, tmp)
	tmp = make([]int, 0)
	for k, _ := range rbt.IterRange(5, 8) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{5, 6, 7}, tmp)

	tmp = make([]int, 0)
	for k, _ := range rbt.IterRange(0, 4) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{1, 2, 3}, tmp)

	tmp = make([]int, 0)
	for k, _ := range rbt.IterRange(12, 99) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, 0, len(tmp))

	tmp = make([]int, 0)
	for k, _ := range rbt.IterBE(5) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{5, 6, 7, 8}, tmp)

	tmp = make([]int, 0)
	for k, _ := range rbt.IterBE(15) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, 0, len(tmp))

	tmp = make([]int, 0)
	for k, _ := range rbt.IterLE(5) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, tmp)

	tmp = make([]int, 0)
	for k, _ := range rbt.IterLE(0) {
		tmp = append(tmp, k)
	}
	assert.Equal(t, 0, len(tmp))

	m := make(map[int]int)
	for k, v := range rbt.Iter() {
		m[k] = v
	}
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8}, m)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, rbt.Keys())
	for i := range 18 {
		if i >= 1 && i <= 8 {
			v, b := rbt.TryGet(i)
			assert.Equal(t, i, v)
			assert.True(t, b)
		} else {
			_, b := rbt.TryGet(i)
			assert.False(t, b)
		}
	}
	for i := range 18 {
		if i >= 1 && i <= 8 {
			v, b := rbt.Delete(i)
			assert.Equal(t, i, v)
			assert.True(t, b)
		} else {
			_, b := rbt.Delete(i)
			assert.False(t, b)
		}
	}
	assert.Equal(t, 0, rbt.length)
}
