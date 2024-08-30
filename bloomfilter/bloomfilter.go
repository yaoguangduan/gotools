package algo

import (
	"bytes"
	"encoding/binary"
	"github.com/spaolacci/murmur3"
	"gotools/bitarray"
	"math"
)

type BloomFilter struct {
	bitset *bitarray.SyncBitArray
	k      int
	m      int
	bs     int
}

func New(insertions uint, fpp float64) *BloomFilter {
	m := int(math.Ceil(-float64(insertions) * math.Log(fpp) / (math.Log(2) * math.Log(2))))
	k := max(1, math.Ceil(math.Log(2)*float64(m)/float64(insertions)))

	bf := &BloomFilter{
		bitset: bitarray.New(m),
		k:      int(k),
		m:      m,
	}
	bf.bs = m
	return bf
}

func NewWithInsertion(insertions uint) *BloomFilter {
	return New(insertions, 0.03)
}

func (bf *BloomFilter) Add(data []byte) {
	h1, h2 := murmur3.Sum128(data)
	var ch = h1
	for i := 0; i < bf.k; i++ {
		bf.bitset.Set(int((ch & math.MaxUint64) % uint64(bf.bs)))
		ch += h2
	}
}
func (bf *BloomFilter) AddString(data string) {
	bf.Add([]byte(data))
}

func (bf *BloomFilter) Contains(data []byte) bool {
	h1, h2 := murmur3.Sum128(data)
	var ch = h1
	for i := 0; i < bf.k; i++ {
		if !bf.bitset.Get(int((ch & math.MaxUint64) % uint64(bf.bs))) {
			return false
		}
		ch += h2
	}
	return true
}

func (bf *BloomFilter) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint64(bf.m))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, uint64(bf.k))
	if err != nil {
		return nil, err
	}
	arr := bf.bitset.Uint64Array()
	err = binary.Write(buf, binary.LittleEndian, uint64(len(arr)))
	if err != nil {
		return nil, err
	}
	for _, v := range arr {
		err = binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (bf *BloomFilter) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	var m uint64
	var k uint64
	var bl uint64
	err := binary.Read(buf, binary.LittleEndian, &m)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &k)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &bl)
	if err != nil {
		return err
	}

	bitsetData := make([]uint64, bl)
	for i := range bl {
		var v uint64
		err = binary.Read(buf, binary.LittleEndian, &v)
		if err != nil {
			return err
		}
		bitsetData[i] = v
	}

	bf.m = int(m)
	bf.k = int(k)
	bf.bitset = bitarray.NewFrom(bitsetData)
	bf.bs = int(m)
	return nil
}
