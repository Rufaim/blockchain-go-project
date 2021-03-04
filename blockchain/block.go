package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"time"

	"google.golang.org/protobuf/proto"

	pb "github.com/Rufaim/blockchain/message"
)

type block struct {
	pb.Block
}

func (b *block) setHash() {
	//TODO: make mining more explicit
	var (
		hash  [sha256.Size]byte
		nonce uint64
	)

	data := prepareData(b, hashTargetBits)
	for nonce < math.MaxUint64 {
		dataWithNonce := appendNonceToData(data, nonce)
		hash = sha256.Sum256(dataWithNonce)

		if isHashValid(hash[:]) {
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
	return b.validateWithNumTargetBits(hashTargetBits)
}

func (b *block) validateWithNumTargetBits(numTargetBits int) bool {
	data := prepareData(b, numTargetBits)
	dataWithNonce := appendNonceToData(data, b.Nonce)
	hash := sha256.Sum256(dataWithNonce)

	if bytes.Compare(hash[:], b.Hash) != 0 {
		return false
	}

	return isHashValid(hash[:])
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
