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

func xorshift32by32(nBits int) *big.Int {
    // the idea here is to calculate the xorshift of the highest 32 bit multiple
    // then shift right to get the required bit ammount
    var nOf32Values int = nBits / 32
    var mod32 int = int(nBits % 32)
    if mod32 != 0 {
        nOf32Values += 1
    }

    ret := big.NewInt(0)
    seed := uint(time.Now().UnixNano())
    seed ^= seed << 13
    seed ^= seed >> 17
    seed ^= seed << 5
    seed = seed << 32
    seed = seed >> 32
    ret.Or(ret, big.NewInt(int64(seed)))
    for i := 1; i != nOf32Values; i++ {
        /* Algorithm "xor" from p. 4 of Marsaglia, "Xorshift RNGs" */
        seed ^= seed << 13
        seed ^= seed >> 17
        seed ^= seed << 5

        ret.Lsh(ret, uint(32*i))
        ret.Or(ret, big.NewInt(int64(seed>>32))) // force 32 bits
    }

    // shifting right to get the required bit ammount
    if mod32 != 0 {
        ret.Rsh(ret, uint(32 - mod32))
    }

    return ret
}

func millerRabin() {

    return
}

func main() {
    fmt.Println("begin")
    var bitSizeArray = []int {40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}
    
    for _, size := range bitSizeArray {
        var avgTime time.Duration = 0
        for i := 0; i != 500; i++ {
            start := time.Now()
            xorshift32by32(size)
            end := time.Now()
            time := end.Sub(start)
            avgTime+=time
        }
        fmt.Printf("xorshift ,%v ,%v \n",size, (avgTime/500))
    }

    for _, size := range bitSizeArray {
        var avgTime time.Duration = 0
        for i := 0; i != 500; i++ {
            start := time.Now()
            blumblumshub(5807, 6287, 32, uint(size))
            end := time.Now()
            time := end.Sub(start)
            avgTime+=time
        }
        fmt.Printf("blum blum shub ,%v ,%v \n",size, (avgTime/500))
    }

    // fmt.Printf("blumblumshub(30000000091, 40000000003, 4882516701, 100): %32b\n", blumblumshub(30000000091, 40000000003, 4882516701, 4096))
    // millerRabin()
    // fmt.Printf("final value xorshift: %v\n", xorshift32by32(4096))
    // fmt.Printf("final value blum: %32b\n",blumblumshub(5807, 6287, 32, 4096))
    return
}
