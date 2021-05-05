package main

import (
	"fmt"
)

const (
	ROUNDS = 16
)

type Blowfish struct {
	P [ROUNDS + 2]uint32
	S [4][256]uint32
}

func (bf *Blowfish) f(x uint32) uint32 {
	a := uint8(x & 0xff)
	x >>= 8
	b := uint8(x & 0xff)
	x >>= 8
	c := uint8(x & 0xff)
	x >>= 8
	d := uint8(x)

	h := bf.S[0][d] + bf.S[1][c]
	h ^= bf.S[2][b]
	return h + bf.S[3][a]
}

func (bf *Blowfish) EncryptBlock(l, r *uint32) {
	L := *l
	R := *r

	for i := 0; i < ROUNDS; i += 2 {
		L ^= bf.P[i]
		R ^= bf.f(L)
		R ^= bf.P[i+1]
		L ^= bf.f(R)
	}

	*l, *r = R^bf.P[ROUNDS+1], L^bf.P[ROUNDS]
}

func (bf *Blowfish) DecryptBlock(l, r *uint32) {
	L := *l
	R := *r

	for i := ROUNDS; i > 0; i -= 2 {
		L ^= bf.P[i+1]
		R ^= bf.f(L)
		R ^= bf.P[i]
		L ^= bf.f(R)
	}

	*l, *r = R^bf.P[0], L^bf.P[1]
}

func NewBlowfish(key []byte) *Blowfish {
	bf := &Blowfish{}

	keyLen := len(key)

	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j++ {
			bf.S[i][j] = INIT_VECTOR[ROUNDS+2+i*256+j]
		}
	}

	for i, k := 0, 0; i < (ROUNDS + 2); i++ {
		data := uint32(0)
		for j := 0; j < 4; j++ {
			data = (data << 8) | uint32(key[k])
			k += 1
			if k >= keyLen {
				k = 0
			}
		}
		bf.P[i] = INIT_VECTOR[i] ^ data
	}

	datal := uint32(0)
	datar := uint32(0)

	for i := 0; i < (ROUNDS + 2); i += 2 {
		bf.EncryptBlock(&datal, &datar)
		bf.P[i] = datal
		bf.P[i+1] = datar
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j += 2 {
			bf.EncryptBlock(&datal, &datar)
			bf.S[i][j] = datal
			bf.S[i][j+1] = datar
		}
	}
	return bf
}

/**
 * https://en.wikipedia.org/wiki/Blowfish_(cipher)
 * https://github.com/piotrpsz/Blowfish
 */
func main() {
	// L := uint32(1)
	L := uint32(0xdf333fd2)
	// R := uint32(2)
	R := uint32(0x30a71bb4)
	bf := NewBlowfish([]byte("TESTKEY"))

	fmt.Println("Initialized!")
	bf.DecryptBlock(&L, &R)
	fmt.Printf("%08x, %08x\n", L, R)
	bf.EncryptBlock(&L, &R)
	fmt.Printf("%08x, %08x\n", L, R)
}
