//+build tests

package blockchain

import (
	"testing"

	"github.com/Rufaim/blockchain/base58"
	pb "github.com/Rufaim/blockchain/message"
)

func TestIsGenesis(t *testing.T) {
	t.Run("Mock", func(t *testing.T) {
		b := block{pb.Block{}}
		if b.IsGenesis() != true {
			t.Errorf("Mock block should be valid as genesis")
		}
	})
	t.Run("Genesis", func(t *testing.T) {
		// Note: this test is valid only for hashTargetBits<=24 and checking firts bytes to be zeros
		b := block{pb.Block{
			Timestamp: 1614397626,
			Hash:      base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")),
			Nonce:     31121925,
		}}
		if b.IsGenesis() != true {
			t.Errorf("Output of genesis block constructor should be valid as genesis")
		}
	})
	t.Run("Not-Genesis", func(t *testing.T) {
		// Note: this test is valid only for hashTargetBits<=24 and checking firts bytes to be zeros
		b := block{pb.Block{
			Timestamp: 1614398030,
			Hash:      base58.Base58Decode([]byte("1119QRRT5j2PpdU9ntYNZt6ZND3c5rTxtZcD9zyEeGA")),
			PrevHash:  base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")),
			Nonce:     3942121,
		}}
		if b.IsGenesis() != false {
			t.Errorf("Output of new block constructor should not be valid as genesis")
		}
	})
}

func TestValidate(t *testing.T) {
	// Note: this test is valid only for and checking firts bytes to be zeros

	const numTargetBits int = 24

	tests := []struct {
		name    string
		b       block
		isValid bool
	}{
		{"Valid Genesis", block{pb.Block{
			Timestamp: 1614397626,
			Hash:      base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")),
			Nonce:     31121925,
		}}, true},
		{"Valid non-genesis", block{pb.Block{
			Timestamp: 1614398030,
			Hash:      base58.Base58Decode([]byte("1119QRRT5j2PpdU9ntYNZt6ZND3c5rTxtZcD9zyEeGA")),
			PrevHash:  base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")),
			Nonce:     3942121,
		}}, true},
		{"Invalid hash Genesis", block{pb.Block{
			Timestamp: 1614397626,
			Hash:      base58.Base58Decode([]byte("1234")),
			Nonce:     31121925,
		}}, false},
		{"Invalid nonce Genesis", block{pb.Block{
			Timestamp: 1614397626,
			Hash:      base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")),
			Nonce:     0,
		}}, false},
		{"invalid previous hash non-genesis", block{pb.Block{
			Timestamp: 1614398030,
			Hash:      base58.Base58Decode([]byte("1119QRRT5j2PpdU9ntYNZt6ZND3c5rTxtZcD9zyEeGA")),
			PrevHash:  base58.Base58Decode([]byte("1234")),
			Nonce:     3942121,
		}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.validateWithNumTargetBits(numTargetBits); got != tt.isValid {
				t.Errorf("Validate() = %t, expected: %t", got, tt.isValid)
			}
		})
	}

}
