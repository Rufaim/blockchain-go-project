package blockchain

import (
	"bytes"
	"crypto/sha256"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
	"google.golang.org/protobuf/proto"
)

func newTxInput(id []byte, outid int, pubkey []byte) *pb.TXInput {
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

func getAllTransactionInputsIds(tx *pb.Transaction) [][]byte {
	refTxIds := make([][]byte, 0, len(tx.Inps))

	for _, in := range tx.Inps {
		refTxIds = append(refTxIds, in.Id)
	}
	return refTxIds
}

//hashTransactions returns a hash of a transaction slice
//it does not verify transactions or its hash
//assuming id is a valid hash been taken during transaction creation
func hashTransactions(txs []*pb.Transaction) []byte {
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

//hashTransactions returns a hash of a transaction
//it is used in transaction creation
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

func trimCopyTransaction(tx *pb.Transaction) (*pb.Transaction, error) {
	txCopy, ok := proto.Clone(tx).(*pb.Transaction)
	if !ok {
		return nil, ErrorTransactionCopyFailed
	}
	for i := range tx.Inps {
		txCopy.Inps[i].Signature = nil
		txCopy.Inps[i].PubKey = nil
	}
	return txCopy, nil
}
