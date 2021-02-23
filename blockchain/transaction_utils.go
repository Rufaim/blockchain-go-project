package blockchain

import (
	"bytes"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
)

func InputUsesKey(input *pb.TXInput, pubKeyHash []byte) bool {
	inputKeyHash := wallet.HashPubKey(input.PubKey)
	return bytes.Compare(inputKeyHash, pubKeyHash) == 0
}

func OutputIsLockedWithKey(output *pb.TXOutput, pubKeyHash []byte) bool {
	return bytes.Compare(output.PubKeyHash, pubKeyHash) == 0
}

func isTransactionCoinbase(tx *pb.Transaction) bool {
	return len(tx.Inps) == 1 && len(tx.Outs) == 0 && tx.Inps[0].OutId == -1
}
