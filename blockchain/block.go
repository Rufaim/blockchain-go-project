package blockchain

import (
	"crypto/sha256"
	"math"
	"math/big"
	"time"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/golang/protobuf/proto"
)

// type block struct {
// 	Timestamp    int64
// 	Transactions []*pb.Transaction
// 	PrevHash     []byte
// 	Hash         []byte
// 	Nonce        uint
// }

type block struct {
	pb.Block
}

func (b *block) setHash() {
	var hashInt big.Int
	var hash [sha256.Size]byte

	var nonce uint64
	for nonce < math.MaxUint64 {
		data := prepareData(b, nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(hashTargetValue) < 0 {
			break
		}
		nonce++
	}
	b.Hash = hash[:]
	b.Nonce = nonce
}

//IsGenesis method returns true if the block is
// the first in chain
func (b *block) IsGenesis() bool {
	return len(b.PrevHash) == 0
}

//Validate returns true if block hash is valid in terms of Proof of Work
func (b *block) Validate() bool {
	var hashInt big.Int

	data := prepareData(b, b.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(hashTargetValue) == -1

	return isValid
}

//Serialize returs a representation of a block as a byte array
func (b *block) Serialize() ([]byte, error) {
	encoded, err := proto.Marshal(b)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func newBlock(txs []*pb.Transaction, prevBlockHash []byte) *block {
	block := &block{pb.Block{
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevBlockHash,
		Transactions: txs,
	}}
	block.setHash()
	return block
}

func newGenesisBlock(coinbase *pb.Transaction) *block {
	genesisBlock := newBlock([]*pb.Transaction{coinbase}, []byte{})
	genesisBlock.setHash()
	return genesisBlock
}
