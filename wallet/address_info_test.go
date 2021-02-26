package wallet_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Rufaim/blockchain/base58"
	"github.com/Rufaim/blockchain/wallet"
)

func TestGetAddressInfo(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    *wallet.AddressInfo
	}{
		{"Address 1", "1AJfyqdw2jwRd9LY6ZBxaRYFtX5koGRmMm", &wallet.AddressInfo{byte(0), destringify("2RUR4KwaAorNnJen1sgV5uxW1Xwn"), destringify("7CQS3M")}},
		{"Address 2", "1ASP8Fy2LMi6BTrqwewnrRVQwvCr6wQHEg", &wallet.AddressInfo{byte(0), destringify("2SekMZdamME3qbjfK9uphEDod3wS"), destringify("7A5Gmn")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wallet.GetAddressInfo([]byte(tt.address)); !compareAddressinfo(got, tt.want) {
				t.Errorf("GetAddressInfo(%s) = %v, want %v", tt.address, got, tt.want)
			}
		})
	}
}

func TestWallet_GetAddressInfo(t *testing.T) {
	for i := 1; i <= 3; i++ {
		t.Run(fmt.Sprintf("Real wallet %d", i), func(t *testing.T) {
			w, _ := wallet.NewWallet()
			w_addr := w.GetAddress()
			w_addr_info := w.GetAddressInfo()
			if !compareAddressinfo(wallet.GetAddressInfo(w_addr), w_addr_info) {
				t.Error("Address info does not match")
			}
		})
	}
}

func compareAddressinfo(a, b *wallet.AddressInfo) bool {
	return a.Version == b.Version && bytes.Compare(a.PublicKeyHash, b.PublicKeyHash) == 0 && bytes.Compare(a.ChechSum, b.ChechSum) == 0
}

func destringify(str string) []byte {
	return base58.Base58Decode([]byte(str))
}
