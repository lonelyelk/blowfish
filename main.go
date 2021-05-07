package main

import (
	"encoding/binary"
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

func (bf *Blowfish) Encrypt(text []byte) []byte {
	blocks := []uint32{uint32(len(text))}
	toAppend := 4 - len(text)%8
	if toAppend < 0 {
		toAppend += 8
	}
	for i := 0; i < toAppend; i++ {
		text = append(text, 0)
	}
	for i := 0; i < len(text); i += 4 {
		blocks = append(blocks, binary.BigEndian.Uint32(text[i:i+4]))
	}
	for i := 0; i < len(blocks); i += 2 {
		bf.EncryptBlock(&blocks[i], &blocks[i+1])
	}
	cryptoText := make([]byte, len(blocks)*4)
	for i, block := range blocks {
		cryptoText[i*4] = byte((block >> 24) & 0xff)
		cryptoText[i*4+1] = byte((block >> 16) & 0xff)
		cryptoText[i*4+2] = byte((block >> 8) & 0xff)
		cryptoText[i*4+3] = byte(block & 0xff)
	}
	return cryptoText
}

func (bf *Blowfish) Decrypt(cryptoText []byte) []byte {
	blocks := []uint32{}
	for i := 0; i < len(cryptoText); i += 4 {
		blocks = append(blocks, binary.BigEndian.Uint32(cryptoText[i:i+4]))
	}
	for i := 0; i < len(blocks); i += 2 {
		bf.DecryptBlock(&blocks[i], &blocks[i+1])
	}
	length := int(blocks[0])
	text := make([]byte, (len(blocks)-1)*4)
	for i, block := range blocks[1:] {
		text[i*4] = byte((block >> 24) & 0xff)
		text[i*4+1] = byte((block >> 16) & 0xff)
		text[i*4+2] = byte((block >> 8) & 0xff)
		text[i*4+3] = byte(block & 0xff)
	}
	return text[:length]
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

func main() {
	bf := NewBlowfish([]byte("Maximum key length is four hundred and eight (448) bits!"))
	text := []byte("CryptoHello CryptoWorld!")
	cryptoText := bf.Encrypt(text)

	for _, b := range cryptoText {
		fmt.Printf("0x%02x ", b)
	}
	fmt.Println()

	fmt.Println(string(bf.Decrypt(cryptoText)))
}
