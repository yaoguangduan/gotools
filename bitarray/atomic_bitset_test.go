package bitarray

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGet(t *testing.T) {
	length := 200
	ab := New(length)
	for i := 0; i < length; i += 2 {
		ab.Set(i)
		assert.True(t, ab.Get(i))
	}
	for i := 1; i < length; i += 2 {
		assert.False(t, ab.Get(i))
	}
	assert.Equal(t, "[1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,1010,0000,0000,0000,0000,0000,0000,0000,0000,0000,0000,0000,0000,0000,0000]", fmt.Sprintf("%s", ab))
	assert.Equal(t, length, ab.Len())
	assert.Equal(t, length/2, ab.BitCnt())
}

func TestPutUint64(t *testing.T) {
	ba := New(64)
	ba.PutUint64(0, 1)
	fmt.Println(ba)
	assert.Equal(t, true, ba.Get(0))
}
