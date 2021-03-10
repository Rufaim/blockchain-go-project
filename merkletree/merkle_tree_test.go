//+build tests

package merkletree_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Rufaim/blockchain/merkletree"
)

func TestMerkleTree(t *testing.T) {
	tests := []struct {
		name  string
		input [][]byte
		want  string
	}{
		{"One", [][]byte{[]byte("123")}, "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"},
		{"Two", [][]byte{[]byte("123"), []byte("234")}, "828a7295a5a3642ff7e9b931e590407c786a4324979f367213a39db89da0fa40"},
		{"Three", [][]byte{[]byte("456"), []byte("345"), []byte("234")}, "f5e5d0e182d9bec396bb9657dc34d1931d4583e474b637c59653a9030009f9fd"},
		{"Seven", [][]byte{[]byte("48798415"), []byte("345"), []byte("878"), []byte("364"), []byte("665"), []byte("7987724"), []byte("81165")}, "1910dd04180102e7de4394b6d5f1ce2d4e705165b43863874e66f2694d98959c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := merkletree.MerkleTreeHash(tt.input)
			sgot := fmt.Sprintf("%x", got)
			if strings.Compare(sgot, tt.want) != 0 {
				t.Errorf("MerkleTreeHash() = %s, want %s", sgot, tt.want)
			}

			got_ := merkletree.MerkleTreeHash(tt.input)
			if bytes.Compare(got, got_) != 0 {
				t.Error("MerkleTreeHash() is stochastic")
			}
		})
	}
}
