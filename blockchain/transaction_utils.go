package blockchain

import (
	"strings"

	pb "github.com/Rufaim/blockchain/message"
)

func inputUsesKey(input *pb.TXInput, pubKey string) bool {
	return strings.Compare(input.PubKey, pubKey) == 0
}

func outputIsLockedWithKey(output *pb.TXOutput, pubKey string) bool {
	return strings.Compare(output.PubKey, pubKey) == 0
}

func isTransactionCoinbase(tx *pb.Transaction) bool {
	return len(tx.Inps) == 1 && len(tx.Outs) == 0 && tx.Inps[0].OutId == -1
}
