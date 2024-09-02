package bitarray

import (
	"fmt"
	"gotools/i64adder"
	"iter"
	"math/bits"
	"strings"
	"sync/atomic"
)

const uint64Bit = 6
const bitPerUnit = 64

type SyncBitArray struct {
	data   []atomic.Uint64
	bitCnt *i64adder.Adder
	len    int
}

func New(size int) *SyncBitArray {
	if size <= 0 {
		panic("size must be greater than zero")
	}
	ul := (size + bitPerUnit - 1) / bitPerUnit
	sba := &SyncBitArray{data: make([]atomic.Uint64, ul), bitCnt: i64adder.New()}
	sba.len = size
	return sba
}

// Len array length
func (ab *SyncBitArray) Len() int {
	return ab.len
}

// Set given index to 1b
func (ab *SyncBitArray) Set(index int) bool {
	if index < 0 || index >= ab.len {
		panic(fmt.Sprintf("index %d out of range %d", index, ab.len))
	}
	aIdx := index >> uint64Bit
	mask := uint64(1) << (index % 64)

	var oldValue uint64
	var newValue uint64
	for {
		oldValue = ab.data[aIdx].Load()
		newValue = oldValue | mask
		if oldValue == newValue {
			return false
		}
		if ab.data[aIdx].CompareAndSwap(oldValue, newValue) {
			ab.bitCnt.Add(1)
			return true
		}
	}
}

// Get return true if given index in bitarray is 1b
func (ab *SyncBitArray) Get(index int) bool {
	if index < 0 || index >= ab.len {
		panic(fmt.Sprintf("index %d out of range %d", index, ab.len))
	}
	return ab.data[index>>uint64Bit].Load()&(1<<(index%64)) != 0
}

// BitCnt bit 1 count
func (ab *SyncBitArray) BitCnt() int {
	return int(ab.bitCnt.Sum())
}

// Iter [index,set flag]
func (ab *SyncBitArray) Iter() iter.Seq2[int, bool] {
	return func(yield func(int, bool) bool) {
		for i := range ab.len {
			if !yield(i, ab.Get(i)) {
				break
			}
		}
	}
}

func (ab *SyncBitArray) Unset(index int) bool {
	if index < 0 || index >= ab.len {
		panic(fmt.Sprintf("index %d out of range %d", index, ab.len))
	}
	aIdx := index >> uint64Bit
	mask := ^(uint64(1) << (index % 64))

	var oldValue uint64
	var newValue uint64
	for {
		oldValue = ab.data[aIdx].Load()
		newValue = oldValue & mask
		if oldValue == newValue {
			return false
		}
		if ab.data[aIdx].CompareAndSwap(oldValue, newValue) {
			ab.bitCnt.Add(-1)
			return true
		}
	}
}

// Clear all bits in the array
func (ab *SyncBitArray) Clear() {
	for i := range ab.data {
		oldValue := ab.data[i].Swap(0)
		ab.bitCnt.Add(-int64(bits.OnesCount64(oldValue)))
	}
}

// PutUint64 put all value bit into uint64 idx
func (ab *SyncBitArray) PutUint64(idx int, value uint64) {
	var update bool
	var old uint64
	var neu uint64
	for {
		old = ab.data[idx>>uint64Bit].Load()
		neu = old | value
		if old == neu {
			break
		}
		if ab.data[idx>>uint64Bit].CompareAndSwap(old, neu) {
			update = true
			break
		}
	}
	if update {
		add := bits.OnesCount64(neu) - bits.OnesCount64(old)
		ab.bitCnt.Add(int64(add))
	}
}

// Uint64Array convert bitarray to uint64 slice
func (ab *SyncBitArray) Uint64Array() []uint64 {
	ret := make([]uint64, len(ab.data))
	for i := range ret {
		ret[i] = ab.data[i].Load()
	}
	return ret
}

// NewFrom new lock free bitarray from exist array
func NewFrom(data []uint64) *SyncBitArray {
	ab := &SyncBitArray{data: make([]atomic.Uint64, len(data)), bitCnt: i64adder.New()}
	for i := range data {
		ab.data[i] = atomic.Uint64{}
		ab.data[i].Store(data[i])
		cnt := bits.OnesCount64(data[i])
		ab.bitCnt.Add(int64(cnt))
	}
	ab.len = len(data) * 64
	return ab
}

func (ab *SyncBitArray) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	u64a := ab.Uint64Array()
	for i, v := range u64a {
		binaryStr := fmt.Sprintf("%064b", v)
		for j, ch := range reverse(binaryStr) {
			if j > 0 && j%4 == 0 {
				sb.WriteRune(',')
			}
			sb.WriteRune(ch)
		}
		if i != len(u64a)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteString("]")
	return sb.String()
}
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
