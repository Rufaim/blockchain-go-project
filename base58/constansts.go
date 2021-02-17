package base58

import "math/big"

var (
	alphabetBTC    = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	decodeAlphabet = make(map[byte]int8)
	baseBig        = big.NewInt(int64(len(alphabetBTC)))
	zeroBig        = big.NewInt(0)
)

func init() {
	for i, b := range alphabetBTC {
		decodeAlphabet[b] = int8(i)
	}
}
