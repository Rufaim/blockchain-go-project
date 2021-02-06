package blockchain

import (
	"bytes"
	"crypto/sha256"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/golang/protobuf/proto"
)

//HashTransactions returns a hash of a transaction slice
func HashTransactions(txs []*pb.Transaction) []byte {
	var (
		txHashes [][]byte
		hash     [sha256.Size]byte
	)

	for _, tx := range txs {
		txHashes = append(txHashes, tx.Id)
	}
	hash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}

func hashTransaction(tx *pb.Transaction) []byte {
	hash := sha256.Sum256(serializeTransaction(tx))
	return hash[:]
}

//Serialize returns a byte version of transaction
func serializeTransaction(tx *pb.Transaction) []byte {
	encoded, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	return encoded
}

func NewTransaction(data []byte) *pb.Transaction {
	hash := sha256.Sum256(data)
	tx := &pb.Transaction{
		Id:   hash[:],
		Data: data,
	}
	return tx
}
