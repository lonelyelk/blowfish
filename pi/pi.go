package pi

import (
	"math"
	"math/big"
)

func modf(f float64) float64 {
	_, frac := math.Modf(f)
	if frac < 0 {
		frac += 1
	}
	return frac
}

func sum(n uint64, j uint64) (result float64) {
	denom := float64(j)
	base := big.NewInt(16)
	for k := float64(0); k <= float64(n); k++ {
		pow := big.NewInt(int64(n) - int64(k))
		den := big.NewInt(int64(denom))
		powermod := new(big.Int)
		powermod.Exp(base, pow, den)
		result = modf(result + float64(powermod.Int64())/denom)
		denom += 8
	}
	num := float64(1.0) / 16
	for frac := num / denom; frac > math.Nextafter(0, 1); frac = num / denom {
		result += frac
		num /= 16
		denom += 8
	}
	return
}

func PiDigit(n uint64) uint {
	frac := modf(4*sum(n, 1) - 2*sum(n, 4) - sum(n, 5) - sum(n, 6))
	return uint(math.Floor(16 * frac))
}

func Nums32(n uint32) (result uint32) {
	for digit := 0; digit < 8; digit++ {
		result += uint32(PiDigit(uint64(n)*8+uint64(digit))) << ((7 - digit) * 4)
	}
	return
}
