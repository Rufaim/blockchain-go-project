package blockchain

import (
	"bytes"
	"crypto/sha256"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
	"github.com/golang/protobuf/proto"
)

func newTxInput(id []byte, outid int, pubkey, signature []byte) *pb.TXInput {
	return &pb.TXInput{
		Id:     id,
		OutId:  int32(outid),
		PubKey: pubkey,
	}
}

func newTxOutput(amount int, address []byte) *pb.TXOutput {
	wi := wallet.GetAddressInfo(address)
	return &pb.TXOutput{
		Amount:     int32(amount),
		PubKeyHash: wi.PublicKeyHash[:],
	}
}

func InputUsesKey(input *pb.TXInput, pubKeyHash []byte) bool {
	inputKeyHash := wallet.HashPubKey(input.PubKey)
	return bytes.Compare(inputKeyHash, pubKeyHash) == 0
}

func OutputIsLockedWithKey(output *pb.TXOutput, pubKeyHash []byte) bool {
	return bytes.Compare(output.PubKeyHash, pubKeyHash) == 0
}

func getAllNonCoinbaseIds(tx *pb.Transaction) [][]byte {
	refTxIds := make([][]byte, 0, len(tx.Inps))

	for _, in := range tx.Inps {
		if in.OutId != -1 {
			refTxIds = append(refTxIds, in.Id)
		}
	}
	return refTxIds
}

func hashTransaction(tx *pb.Transaction) []byte {
	hash := sha256.Sum256(serializeTransaction(tx))
	return hash[:]
}

func serializeTransaction(tx *pb.Transaction) []byte {
	encoded, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	return encoded
}

func isTransactionCoinbase(tx *pb.Transaction) bool {
	return len(tx.Inps) == 1 && len(tx.Inps[0].Id) == 0 && tx.Inps[0].OutId == -1
}
