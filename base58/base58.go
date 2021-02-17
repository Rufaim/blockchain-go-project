package base58

import (
	"bytes"
	"math/big"
	"strings"
)

func Base58Encode(input []byte) []byte {
	var result strings.Builder

	x := new(big.Int).SetBytes(input)
	mod := &big.Int{}

	for x.Cmp(zeroBig) != 0 {
		x.DivMod(x, baseBig, mod)
		result.WriteByte(alphabetBTC[mod.Int64()])
	}

	for _, b := range input {
		if b != 0x00 {
			break
		}
		result.WriteByte(alphabetBTC[0])
	}

	// reversing the array
	reversed := []byte(result.String())
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	base := big.NewInt(1)
	zeros := 0

	for i, b := range input {
		if b != alphabetBTC[0] {
			zeros = i
			break
		}
	}

	for i := len(input) - 1; i >= zeros; i-- {
		idx := decodeAlphabet[input[i]]
		tmp := big.NewInt(0)
		tmp.Mul(base, big.NewInt(int64(idx)))
		result.Add(result, tmp)
		base.Mul(base, baseBig)
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeros), decoded...)

	return decoded
}
