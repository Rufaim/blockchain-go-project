package wallet

import (
	"strings"
	"testing"

	"github.com/Rufaim/blockchain/base58"
	pb "github.com/Rufaim/blockchain/message"
	"google.golang.org/protobuf/proto"
)

func TestWallet_GetAddress(t *testing.T) {
	tests := []struct {
		name    string
		pubkey  []byte
		address string
	}{
		{"Mock 1", []byte("1234"), getVersionString() + "Q7GAua66sVqEpHogKWbmLaLC7NYtvh5p7"},
		{"Mock 2", []byte("PsychicTandemWarElephant"), getVersionString() + "3zmG2Mgu2HzNxbEL2Y2T3v4Fxcwvrabc1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock_w := &Wallet{
				PublicKey: tt.pubkey,
			}
			if got := string(mock_w.GetAddress()); strings.Compare(got, tt.address) != 0 {
				t.Errorf("Wallet.GetAddress() = %v, want %v", got, tt.address)
			}
		})
	}

	t.Run("Real wallet", func(t *testing.T) {
		str := "2VeSDJUmhujKZNtA8L7YbgfkoFkrBSysCDo9PXhCUG36PRkfkEdepvhqTWrgiXg7vqWvDQzLNAPA62mLRkkh2BZW1TtcUqnTxpDKArrjQbCNxoouLrKwfpp8bSpGdaioKNXDwQ4EkRS1fKjWr8H5EV5WXx18EFeFpYUvDLtj7ngNGfgET9Hi1j98CJnjXsxj35pWjg42eQBz2hGAgP8Dr36thqcR3Zq7NG1pSTNPa1saBtBggFAQzr46hJxUghZKDfgUKEu2WPCYueNJYA9kvcGo7uBQqgff1cfPRmZas3agNPvkw1GxXn9ZSinzUJCUxGN5cVGAf5FtEJVG2aAw6mJaQBp1rT5AzbheP3uZFzd8nBT8Jok94jPcgXxvfqcvtzSm6cNfT7PdMV3B4nusctLitB3ftzPU9hDNbYer8YVcDM9qeExH5q4KmqFTwEUL75vPAperrzMH6xfUvzgAK92G44"
		address := getVersionString() + "BVrKfQX7NzqebGE5D6zmMr7yXbSkSyYzE"
		pw := &pb.Wallet{}
		if err := proto.Unmarshal(base58.Base58Decode([]byte(str)), pw); err != nil {
			panic(err)
		}
		w := NewFromProto(pw)
		if got := string(w.GetAddress()); strings.Compare(got, address) != 0 {
			t.Errorf("Address error expected: %s ; got: %s", address, got)
		}
	})
}

func getVersionString() string {
	return string(byte('1') + Version)
}
