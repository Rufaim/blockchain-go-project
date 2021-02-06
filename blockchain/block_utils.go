package blockchain

import (
	"bytes"
	"strconv"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/golang/protobuf/proto"
)

func prepareData(b *block, nonce uint64) []byte {
	blockDataBytes := [][]byte{
		b.PrevHash,
		HashTransactions(b.Transactions),
		intToBytes(b.Timestamp, 16),
		intToBytes(int64(hashTargetBits), 16),
		uintToBytes(nonce, 16),
	}

	return bytes.Join(blockDataBytes, []byte{})
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
