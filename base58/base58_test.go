//+build tests

package base58_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/Rufaim/blockchain/base58"
)

func getBytestFromStringInt(in string) []byte {
	i, _ := new(big.Int).SetString(in, 10)
	return bytes.Join([][]byte{[]byte{byte(0x00)}, i.Bytes()}, []byte{})
}

func TestBase58Encoder(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		output []byte
	}{
		{"Simple", []byte("12345"), []byte("6YvUFcg")},
		{"Complex", []byte("Ax45rxbetbs13e124rdslkvnwosdjv2756aveaf9ewmrr03mf9emz-a"), []byte("QJJnjztRDbgGytied4gp4Pt8US7jqST2xLC7yNpLD62MAE4kXHz1LZ7jEDkPERai8Uar3kWUUji")},
		{"Real data", getBytestFromStringInt("3289723086273072581179791804083638811218952016296592203606"), []byte("1DEQEZtucA3aS1wuRG2eHyXPnQHkA1ERaM")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if encoded := base58.Base58Encode(tt.input); bytes.Compare(encoded, tt.output) != 0 {
				t.Errorf("Base58Encoder() expected: %s, got: %s", string(tt.output), string(encoded))
			}
		})
	}
}

func TestBase58(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"Simple", []byte("12345")},
		{"Complex", []byte("Ax45rxbetbs13e124rdslkvnwosdjv2756aveaf9ewmrr03mf9emz-a")},
		{"Real data", getBytestFromStringInt("3289723086273072581179791804083638811218952016296592203606")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := base58.Base58Encode(tt.input)
			decoded := base58.Base58Decode(encoded)

			if bytes.Compare(tt.input, decoded) != 0 {
				t.Errorf("input sequence: %s, decoded: %s", string(tt.input), string(decoded))
			}
		})
	}
}
