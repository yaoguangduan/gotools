package algo

import (
	"cmp"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"maps"
	"os"
	"runtime/pprof"
	"slices"
	"testing"
	"time"
)

type NoOrder struct {
	in int
}

func TestOp(t *testing.T) {
	sl := NewSkipList[string, string]()
	sl.Put("zook", "nonono")
	sl.Put("join", "hhh")
	sl.Put("alen", "good")
	assert.Equal(t, "good", sl.Get("alen"))
	assert.Equal(t, "hhh", sl.Get("join"))
	assert.Equal(t, "nonono", sl.Get("zook"))
	assert.Equal(t, "", sl.Get("nonono"))
	assert.False(t, sl.Has("not exist"))
	assert.Equal(t, "def val", sl.GetOr("nothing", "def val"))
	sl.Put("alen", "new val")
	assert.Equal(t, "new val", sl.Get("alen"))
	assert.True(t, sl.Has("alen"))
	assert.True(t, sl.Has("join"))
	assert.True(t, sl.Has("zook"))
	assert.False(t, sl.Has("nonono"))
	assert.Equal(t, 3, sl.Len())
	for k := range sl.Iter() {
		t.Log(k)
	}
	list := slices.Collect(sl.Iter())
	assert.Equal(t, 3, len(list))
	assert.True(t, slices.Contains(list, "alen"))
	assert.Equal(t, []string{"alen", "join", "zook"}, list)
	assert.Equal(t, "hhh", sl.Delete("join"))
	assert.Equal(t, "", sl.Delete("nonono"))
	assert.Equal(t, 2, sl.Len())
	list = slices.Collect(sl.Iter())
	assert.Equal(t, []string{"alen", "zook"}, list)
	sl.Clear()
	assert.Equal(t, "", sl.Get("alen"))
	assert.Equal(t, "", sl.Get("join"))
	assert.Equal(t, "", sl.Get("zook"))
	assert.Equal(t, 0, sl.Len())
	list = slices.Collect(sl.Iter())
	assert.Nil(t, list)
	assert.Equal(t, 0, len(list))
	for k := range sl.Iter() {
		t.Log(k)
	}

	slf := NewSkipListWithCmp[NoOrder, int](func(a, b NoOrder) int {
		return cmp.Compare(a.in, b.in)
	})
	slf.Put(NoOrder{in: 9}, 1)
	slf.Put(NoOrder{in: 2}, 8)
	slf.Put(NoOrder{in: 6}, 4)
	slf.Put(NoOrder{in: 3}, 7)
	slf.Put(NoOrder{in: 7}, 3)
	slf.Put(NoOrder{in: 4}, 6)
	slf.Put(NoOrder{in: 1}, 9)
	slf.Put(NoOrder{in: 8}, 2)
	slf.Put(NoOrder{in: 5}, 5)
	var i = 1
	var j = 9
	for k, v := range slf.Iter2() {
		assert.Equal(t, i, k.in)
		assert.Equal(t, j, v)
		i++
		j--
	}
	rg := maps.Collect(slf.IterRange(NoOrder{3}, NoOrder{5}))
	assert.Equal(t, 2, len(rg))
	assert.Equal(t, map[NoOrder]int{NoOrder{3}: 7, NoOrder{4}: 6}, rg)
	log.Println(rg)
	empty := maps.Collect(slf.IterRange(NoOrder{5}, NoOrder{3}))
	assert.Equal(t, 0, len(empty))
	log.Println(empty)

	be := maps.Collect(slf.IterBE(NoOrder{7}))
	assert.Equal(t, 3, len(be))
	assert.Equal(t, map[NoOrder]int{NoOrder{7}: 3, NoOrder{8}: 2, NoOrder{9}: 1}, be)
	t.Log(be)
	empty = maps.Collect(slf.IterBE(NoOrder{10}))
	assert.Equal(t, 0, len(empty))
	log.Println(empty)

	le := maps.Collect(slf.IterLE(NoOrder{3}))
	assert.Equal(t, 3, len(le))
	assert.Equal(t, map[NoOrder]int{NoOrder{1}: 9, NoOrder{2}: 8, NoOrder{3}: 7}, le)
	t.Log(le)
	empty = maps.Collect(slf.IterLE(NoOrder{0}))
	assert.Equal(t, 0, len(empty))
	log.Println(empty)

}
func TestRepeat(t *testing.T) {
	// 创建文件以保存 CPU 性能数据
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
		return
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f) // 确保在程序结束时关闭文件

	// 启动 CPU 性能分析
	if err = pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile: ", err)
		return
	}
	defer pprof.StopCPUProfile() // 确保在程序结束时停止性能分析
	sl := NewSkipList[int, int]()
	now := time.Now().UnixMilli()
	for i := 0; i < 1000000; i++ {
		sl.Put(i, i)
	}
	fmt.Println("cost", time.Now().UnixMilli()-now)
}
