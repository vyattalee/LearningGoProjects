package utils

import (
	"fmt"
)

import (
	"bytes"

	"strings"
)

const (
	IsSet   int = 1
	IsUnSet int = 0
)

// BitMap represents a bit array.
type BitMap struct {
	Size int
	data []byte
}

// NewBitMap returns a size-length BitMap pointer.
func NewBitMap(size int) *BitMap {
	div, mod := size>>3, size&0x07
	if mod > 0 {
		div++
	}
	return &BitMap{size, make([]byte, div)}
}

// NewBitMapFrom returns a new copyed BitMap pointer which
// NewBitMap.data = other.data[:size].
func NewBitMapFrom(other *BitMap, size int) *BitMap {
	bitmap := NewBitMap(size)

	if size > other.Size {
		size = other.Size
	}

	div := size >> 3

	for i := 0; i < div; i++ {
		bitmap.data[i] = other.data[i]
	}

	for i := div << 3; i < size; i++ {
		if other.Bit(i) == 1 {
			bitmap.Set(i)
		}
	}

	return bitmap
}

// NewBitMapFromBytes returns a BitMap pointer created from a byte array.
func NewBitMapFromBytes(data []byte) *BitMap {
	bitmap := NewBitMap(len(data) << 3)
	copy(bitmap.data, data)
	return bitmap
}

// NewBitMapFromString returns a BitMap pointer created from a string.
func NewBitMapFromString(data string) *BitMap {
	return NewBitMapFromBytes([]byte(data))
}

// Bit returns the bit at index.
func (bitmap *BitMap) Bit(index int) int {
	if index >= bitmap.Size {
		panic("index out of range")
	}

	div, mod := index>>3, index&0x07
	return int((uint(bitmap.data[div]) & (1 << uint(7-mod))) >> uint(7-mod))
}

// set sets the bit at index `index`. If bit is true, set 1, otherwise set 0.
func (bitmap *BitMap) set(index int, bit int) {
	if index >= bitmap.Size {
		panic("index out of range")
	}

	div, mod := index>>3, index&0x07
	shift := byte(1 << uint(7-mod))

	bitmap.data[div] &= ^shift
	if bit > 0 {
		bitmap.data[div] |= shift
	}
}

// Set sets the bit at idnex to 1.
func (bitmap *BitMap) Set(index int) {
	bitmap.set(index, 1)
}

// Unset sets the bit at idnex to 0.
func (bitmap *BitMap) Unset(index int) {
	bitmap.set(index, 0)
}

// Compare compares the prefixLen-prefix of two BitMap.
//   - If BitMap.data[:prefixLen] < other.data[:prefixLen], return -1.
//   - If BitMap.data[:prefixLen] > other.data[:prefixLen], return 1.
//   - Otherwise return 0.
func (bitmap *BitMap) Compare(other *BitMap, prefixLen int) int {
	if prefixLen > bitmap.Size || prefixLen > other.Size {
		panic("index out of range")
	}

	div, mod := prefixLen>>3, prefixLen&0x07
	res := bytes.Compare(bitmap.data[:div], other.data[:div])
	if res != 0 {
		return res
	}

	for i := div << 3; i < (div<<3)+mod; i++ {
		bit1, bit2 := bitmap.Bit(i), other.Bit(i)
		if bit1 > bit2 {
			return 1
		} else if bit1 < bit2 {
			return -1
		}
	}

	return 0
}

// Xor returns the xor value of two BitMap.
func (bitmap *BitMap) Xor(other *BitMap) *BitMap {
	if bitmap.Size != other.Size {
		panic("size not the same")
	}

	distance := NewBitMap(bitmap.Size)
	xor(distance.data, bitmap.data, other.data)

	return distance
}

// String returns the bit sequence string of the BitMap.
func (bitmap *BitMap) String() string {
	div, mod := bitmap.Size>>3, bitmap.Size&0x07
	buff := make([]string, div+mod)

	for i := 0; i < div; i++ {
		buff[i] = fmt.Sprintf("%08b", bitmap.data[i])
	}

	for i := div; i < div+mod; i++ {
		buff[i] = fmt.Sprintf("%1b", bitmap.Bit(div*8+(i-div)))
	}

	return strings.Join(buff, "")
}

// RawString returns the string value of BitMap.data.
func (bitmap *BitMap) RawString() string {
	return string(bitmap.data)
}

// String returns the bit sequence string of the BitMap.
func (bitmap *BitMap) Bytes() []byte {
	return bitmap.data
}
