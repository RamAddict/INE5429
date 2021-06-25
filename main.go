package main

import (
	"fmt"
	"time"
)

/*
// @param
	p and q are primes that should both be congruent to 3 (mod 4)
*/
func blumblumshub(p uint64, q uint64, seed uint64, nBits int) []uint32 {
	var nOf32Values int = nBits / 32
	var mod32 int = int(nBits % 32)
	if mod32 != 0 {
		nOf32Values += 1
	}

	m := p * q

	var nBitArray = make([]uint32, nOf32Values)
	nBitArray[0] |= uint32(seed >> 31)
	xn := seed
	for i := 0; i != nBits; i++ {
		xnPlus1 := uint64(uint64(xn*xn) % m)
		xn = xnPlus1
		if xnPlus1%2 == 0 {
			nBitArray[i/32] |= (uint32(0) << (i % 32))
		} else {
			nBitArray[i/32] |= (uint32(1) << (i % 32))
		}
	}
	return nBitArray
}

func xorshift32by32(nBits int) []uint32 {
	// the idea here is to calculate the xorshift of the highest 32 bit multiple
	// then shift right to get the required bit ammount
	var nOf32Values int = nBits / 32
	var mod32 int = int(nBits % 32)
	if mod32 != 0 {
		nOf32Values += 1
	}

	var nBitArray = make([]uint32, nOf32Values)

	for i := range nBitArray {
		seed := uint32(time.Now().UnixNano())
		/* Algorithm "xor" from p. 4 of Marsaglia, "Xorshift RNGs" */
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5

		nBitArray[i] = seed
	}
	// shifting right to get the required bit ammount
	if mod32 != 0 {
		nBitArray[len(nBitArray)-1] = nBitArray[len(nBitArray)-1] >> (32 - mod32)
	}

	return nBitArray
}

func main() {
	// var bitSizeArray = []int {40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}
	fmt.Println("begin")
	fmt.Println(xorshift32by32(8))
	fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", blumblumshub(5807, 6287, 32, 4096))
	fmt.Println(blumblumshub(5807, 6287, 32, 56))
}
