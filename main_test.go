package main

import (
	"testing"
	"testing/quick"
)

func TestBlowfishBlocks(t *testing.T) {
	symBlock := func(key []byte, l uint32, r uint32) bool {
		if len(key) < 32/8 {
			key = append(key, []byte("test")...)
		}
		if len(key) > 448/8 {
			key = key[0 : 448/8]
		}
		copyL, copyR := l, r
		bf := NewBlowfish(key)
		bf.EncryptBlock(&l, &r)
		bf.DecryptBlock(&l, &r)
		return l == copyL && r == copyR
	}
	if err := quick.Check(symBlock, nil); err != nil {
		t.Error((err))
	}
}
