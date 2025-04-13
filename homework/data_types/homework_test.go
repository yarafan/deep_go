package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type uints interface {
	~uint16 | ~uint32 | ~uint64
}

func ToLittleEndianGen[T uints](number T) T {
	var result T
	var size = int(unsafe.Sizeof(number))
	var bytePos = size - 1

	for i := 0; i < size; i += 1 {
		mask := T(0xFF) << (i * 8)
		result |= ((number & mask) >> (i * 8)) << (bytePos * 8)
		bytePos -= 1
	}

	return result
}

func ToLittleEndian(number uint32) uint32 {
	return (number>>24)&0xFF | (number>>8)&0xFF00 | (number<<8)&0xFF0000 | (number<<24)&0xFF000000
}

func TestСonversionUint32(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)

			result2 := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result2)
		})
	}
}

func TestСonversionUint16(t *testing.T) {
	tests := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x0102,
			result: 0x0201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result2 := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result2)
		})
	}
}

func TestСonversionUint64(t *testing.T) {
	tests := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #3": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result2 := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result2)
		})
	}
}
