package pi

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestPiDigit(t *testing.T) {
	piHexDigitStr := "0x3."
	for i := 0; i < 20; i++ {
		piHexDigitStr += strconv.FormatUint(uint64(PiDigit(uint64(i))), 16)
	}
	piHexDigitStr += "p+0"
	num, err := strconv.ParseFloat(piHexDigitStr, 64)
	if err != nil {
		t.Error(err)
	} else if math.Pi != num {
		t.Errorf("Expected %.20f to equal %.20f\n", num, math.Pi)
	}
}

func TestNums32(t *testing.T) {
	piHexDigitStr := "0x3."
	for i := 0; i < 2; i++ {
		piHexDigitStr += fmt.Sprintf("%08x", Nums32(uint32(i)))
	}
	piHexDigitStr += "p+0"
	num, err := strconv.ParseFloat(piHexDigitStr, 64)
	if err != nil {
		t.Error(err)
	} else if math.Pi != num {
		t.Errorf("Expected %.20f to equal %.20f\n", num, math.Pi)
	}
}
