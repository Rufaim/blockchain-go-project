package wallet

import (
	"fmt"
	"testing"
)

func TestChecksumLenSize(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"Short", []byte("ytoqnc")},
		{"Long", []byte("12345678912356789qwertyuiop[]asdfghjkl;'zxcvbnm,./qwertyuiop[]sdfghjkl;'zxcvbnm,./159753")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksum(tt.input); len(got) != AddressChecksumLen {
				t.Errorf("len of returned chechsum in not equal to %d", AddressChecksumLen)
			}
		})
	}
}

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{"Valid 1", []byte("1ASP8Fy2LMi6BTrqwewnrRVQwvCr6wQHEg"), true},
		{"Valid 2", []byte("1ue22zcdtmBg5Wcpe9rdvp8G6Yz4VTFjC"), true},
		{"Short", []byte("1Wh4bh"), false},
		{"Wrong version", []byte(string(byte('1')+Version+byte(1)) + "ASP8Fy2LMi6BTrqwewnrRVQwvCr6wQHEg"), false},
		{"Invalid checksum", []byte("1EKdwAAQDzwaYBCiZSYcSw6cYqdsmfhzuo"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.input); got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}

	for i := 1; i <= 3; i++ {
		t.Run(fmt.Sprintf("Real wallet %d", i), func(t *testing.T) {
			w, _ := NewWallet()
			if IsValidAddress(w.GetAddress()) != true {
				t.Error("IsValidAddress is failed on a real wallet")
			}
		})
	}
}
