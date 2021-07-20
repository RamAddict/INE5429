package main

import (
	"fmt"
	"math/big"
	"time"
)

/*
// Param
    pa and qa are primes that should both be congruent to 3 (mod 4)
	seed is any random number such that mdc(p*q, seed) = 1
	nBits is the size in bits of the output number
// Return
	the random number
*/
func blumblumshub(pa int64, qa int64, 
				  seed int64, nBits uint) *big.Int {

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
	// set left most bit to 1
	mask := big.NewInt(1).Lsh(big.NewInt(1), uint(nBits))
	ret.Or(ret, mask)
	return ret
}
/*
// Param
	nBits is the size in bits of the output number
// Return
	the random number
*/
func xorshift32by32(nBits int) *big.Int {
	// the idea here is to calculate the xorshift of the highest 32 bit multiple
	// then shift right to get the required bit amount
	var nOf32Values int = nBits / 32
	var mod32 int = int(nBits % 32)
	if mod32 != 0 {
		nOf32Values += 1
	}
 
	ret := big.NewInt(0)
	// seed is hard coded as the time
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
 
	// shifting right to get the required bit amount
	if mod32 != 0 {
		ret.Rsh(ret, uint(32-mod32))
	}
	// set left most bit to 1
	mask := big.NewInt(1).Lsh(big.NewInt(1), uint(nBits))
	ret.Or(ret, mask)
	return ret.Add(ret, big.NewInt(1))
}

func lt(a *big.Int, b *big.Int) bool {
	if a.Cmp(b) == -1 {
		return true
	}
	return false

	// -1 if x < y
	// 0 if x == y
	// +1 if x > y
}

func eq(a *big.Int, b *big.Int) bool {
	if a.Cmp(b) == 0 {
		return true
	}
	return false

	// -1 if x < y
	// 0 if x == y
	// +1 if x > y
}
/*
// Param
	nPossiblePrime is the number that will be tested
	attemps is the number of certainty iterations
// Return
	if the number is prime
*/
func  millerRabin(nPossiblePrime *big.Int, attempts uint) bool {

	// if less than 1 
	if lt(nPossiblePrime, big.NewInt(1)) {
        return false
    }
	
	// if even, stop
    if eq(big.NewInt(0).Mod(nPossiblePrime, big.NewInt(2)), big.NewInt(0)) {
        return false
    }
	// obviously prime
	if eq(nPossiblePrime, big.NewInt(1)) || eq(nPossiblePrime, big.NewInt(2)) || eq(nPossiblePrime, big.NewInt(3)) || eq(nPossiblePrime, big.NewInt(5)){
        return true
    }
	// n -1
	nMinus1 := big.NewInt(0).Sub(nPossiblePrime, big.NewInt(1))
	
	// FORCE COPY
	t := big.NewInt(0).Add(nMinus1, big.NewInt(0))
	s := 0
	for {
		// if odd, stop
		if !eq(big.NewInt(0).Mod(t, big.NewInt(2)), big.NewInt(0)) {
			break
		}
		// increment max exp that divides 
		s++
		// divide by 2
		t.Rsh(t, 1)
	}

	for i := uint(0); i != attempts; i++ {
		a := xorshift32by32(nPossiblePrime.BitLen()/2)
		x := big.NewInt(0).Exp(a, t, nPossiblePrime)
		if eq(x, nMinus1) || eq(x, big.NewInt(1)) {
			// result is good, continue
			continue
		}
		didStop := false
		for i := 0; i < (s-1); i++ {
			x.Exp(x, big.NewInt(2), nPossiblePrime)
			if eq(x, nMinus1) {
				// result is good, continue
				didStop = true
				break
			}
		}
		if !didStop {
			return false
		}
	}
	
    return true
}

/*
// Param
	nPossiblePrime is the number that will be tested
	attemps is the number of certainty iterations
// Return
	if the number is prime
*/
func fermatTest(possiblePrime *big.Int, attempts uint) bool {

    pMinus1 := big.NewInt(0).Sub(possiblePrime, big.NewInt(1))

	for i := uint(0); i != attempts; i++ {
		// generate a random number which is less than p
		a := xorshift32by32(20)
		for {
			if lt(a, possiblePrime) {
				break
			} else {
				a = xorshift32by32(20)
			}
		}
		gcd := big.NewInt(0)
		gcd.GCD(nil, nil, a, possiblePrime)
		if !eq(gcd, big.NewInt(1)) {
			// definately not prime
			return false
		} else {
            isComposite := big.NewInt(0).Exp(a, pMinus1, possiblePrime)
            if !eq(isComposite, big.NewInt(1)) {
                // not prime
                return false
            }
        }
	}
	return true
}

func main() {
	fmt.Println("begin")
	for i := uint(0); i != 5; i++ {
		potentialPrime := xorshift32by32(256)
		if fermatTest(potentialPrime, 20) {
            fmt.Printf("fermatTest says %v  is prime!\n", potentialPrime)
        }
        fmt.Printf("i: %v\n", i)
        a := xorshift32by32(256)
        if millerRabin(a, 20) {
            fmt.Printf("Miller Rabin says %v is prime\n", a)
        }
    }
    

	var bitSizeArray = []int {40, 56, 80, 128, 168, 224, 256, 512, 1024, 2048, 4096}

    fmt.Println("algo, bits, avg_time")

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

	bitSizeArray = []int {40, 56, 80, 128, 168, 224, 256, 512, 1024}
	
	for _, size := range bitSizeArray {
	    var avgTime time.Duration = 0
	    for i := 0; i != 5; i++ {
	        potentialPrime := xorshift32by32(size)
			start := time.Now()
	        if millerRabin(potentialPrime, 20) {
				fmt.Printf("millerRabin says %v  is prime!\n", potentialPrime)
			}
	        end := time.Now()
	        time := end.Sub(start)
	        avgTime+=time
	    }
	    fmt.Printf("MillerRabin ,%v ,%v \n",size, (avgTime/1))
	}

	for _, size := range bitSizeArray {
	    var avgTime time.Duration = 0
	    for i := 0; i != 5; i++ {
	        potentialPrime := xorshift32by32(size)
			start := time.Now()
	        if fermatTest(potentialPrime, 20) {
				fmt.Printf("fermatTest says %v  is prime!\n", potentialPrime)
			}
	        end := time.Now()
	        time := end.Sub(start)
	        avgTime+=time
	    }
	    fmt.Printf("Fermat ,%v ,%v \n",size, (avgTime/1))
	}

	// fmt.Printf("blumblumshub(30000000091, 40000000003, 4882516701, 100): %32b\n", blumblumshub(30000000091, 40000000003, 4882516701, 4096))
	// millerRabin()
	// fmt.Printf("final value xorshift: %v\n", xorshift32by32(4096))
	// fmt.Printf("final value blum: %32b\n",blumblumshub(5807, 6287, 32, 4096))
	return
}
