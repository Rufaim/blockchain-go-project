package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"strconv"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/golang/protobuf/proto"
)

func prepareData(b *block, numTargetBits int) []byte {
	blockDataBytes := [][]byte{
		b.PrevHash,
		HashTransactions(b.Transactions),
		intToBytes(b.Timestamp, 16),
		intToBytes(int64(numTargetBits), 16),
	}

	return bytes.Join(blockDataBytes, []byte{})
}

func appendNonceToData(data []byte, nonce uint64) []byte {
	return bytes.Join([][]byte{data, uintToBytes(nonce, 16)}, []byte{})
}

func isHashValid(hash []byte) bool {
	if len(hash) != sha256.Size {
		return false
	}
	h := new(big.Int).SetBytes(hash)
	return h.Cmp(hashTargetValue) <= 0
}

func intToBytes(n int64, base int) []byte {
	return []byte(strconv.FormatInt(n, base))
}

func uintToBytes(n uint64, base int) []byte {
	return []byte(strconv.FormatUint(n, base))
}

//DeserializeBlock returns a block deserialized from a byte array
func DeserializeBlock(bytesArray []byte) (*block, error) {
	b := &pb.Block{}
	if err := proto.Unmarshal(bytesArray, b); err != nil {
		return nil, err
	}
	return &block{*b}, nil
}
