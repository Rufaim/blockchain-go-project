//+build tests

package blockchain

import (
	"bytes"
	"testing"

	"github.com/Rufaim/blockchain/base58"
)

func TestIntToBytes(t *testing.T) {
	type args struct {
		n    int64
		base int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"Positive_10", args{42, 10}, []byte("42")},
		{"Positive_16", args{10, 16}, []byte("a")},
		{"Negative_10", args{-21, 10}, []byte("-21")},
		{"Negative_16", args{-10, 16}, []byte("-a")},
		{"Zero", args{0, 10}, []byte("0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intToBytes(tt.args.n, tt.args.base); bytes.Compare(got, tt.want) != 0 {
				t.Errorf("intToBytes(%d, %d) = %v, want %v", tt.args.n, tt.args.base, string(got), string(tt.want))
			}
		})
	}
}

func TestUintToBytes(t *testing.T) {
	type args struct {
		n    uint64
		base int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"Positive_10", args{42, 10}, []byte("42")},
		{"Positive_16", args{10, 16}, []byte("a")},
		{"Zero", args{0, 10}, []byte("0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uintToBytes(tt.args.n, tt.args.base); bytes.Compare(got, tt.want) != 0 {
				t.Errorf("intToBytes(%d, %d) = %v, want %v", tt.args.n, tt.args.base, string(got), string(tt.want))
			}
		})
	}
}

func TestIsHashValid(t *testing.T) {
	// Note: this test is valid only for hashTargetBits<=24 and checking firts bytes to be zeros
	tests := []struct {
		name    string
		hash    []byte
		isValid bool
	}{
		{"Valid 1", base58.Base58Decode([]byte("111BDuqkq9Su5hGxsgAth4kKtcXipeLrwLU2fAGFw5A")), true},
		{"Valid 2", base58.Base58Decode([]byte("1119QRRT5j2PpdU9ntYNZt6ZND3c5rTxtZcD9zyEeGA")), true},
		{"Short", base58.Base58Decode([]byte("jaeg734")), false},
		{"Invalid", base58.Base58Decode([]byte("2229QRRT5j2Pp7pbto1Zt6ZND3c5rTxtZcD58tB91z")), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHashValid(tt.hash); got != tt.isValid {
				t.Errorf("isHashValid(%x) = %t, exptected: %t", tt.hash, got, tt.isValid)
			}
		})
	}
}
