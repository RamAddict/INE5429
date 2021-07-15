package main

import (
	"fmt"
	"math/big"
	"time"
)

/*
// @param
	p and q are primes that should both be congruent to 3 (mod 4)
*/
func blumblumshub(pa int64, qa int64, seed int64, nBits uint) *big.Int {

	p := big.NewInt(pa)
	q := big.NewInt(qa)
	m := p.Mul(p, q)

	ret := big.NewInt(0)
	xn := big.NewInt(seed)
	for i := uint(0); i != nBits; i++ {
		xn.Set(xn.Mul(xn, xn).Mod(xn, m))
		mod2 := big.NewInt(1)
		mod2.Mod(xn, big.NewInt(2))

		ret = mod2.Or(ret, (mod2.Lsh(mod2, i)))

	}
	return ret
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

func millerRabin() {

	// trash := big.NewInt(1)

	// a := big.NewInt(1822323232323232328)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", a)
	// b := big.NewInt(10)
	// c := big.NewInt(5)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", b)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", c)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", trash)

	// trash.Mul(b,c)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", a)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", b)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", c)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", trash)
	// trash.Mul(trash,c)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", a)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", b)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", c)
	// fmt.Printf("blumblumshub(11, 23, 3, 56): %32b\n", trash)

	return
}

func main() {
	// var bitSizeArray = []int {40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}
	fmt.Println("begin")
	fmt.Printf("blumblumshub(30000000091, 40000000003, 4882516701, 100): %32b\n", blumblumshub(30000000091, 40000000003, 4882516701, 4096))
	// millerRabin()
	// fmt.Println(xorshift32by32(8))
	// fmt.Println(blumblumshub(5807, 6287, 32, 56))
	return
}
