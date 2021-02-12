package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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

func newTxInput(id []byte, outid int, pubkey string) *pb.TXInput {
	return &pb.TXInput{
		Id:     id,
		OutId:  int32(outid),
		PubKey: pubkey,
	}
}

func newTxOutput(amount int, pubkey string) *pb.TXOutput {
	return &pb.TXOutput{
		Amount: int32(amount),
		PubKey: pubkey,
	}
}

func newTransaction(inps []*pb.TXInput, outs []*pb.TXOutput) *pb.Transaction {
	tx := &pb.Transaction{
		Inps: inps,
		Outs: outs,
	}

	encodedData, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(encodedData)

	tx.Id = hash[:]
	return tx
}

func NewCoinbaseTX(to string) *pb.Transaction {
	data := fmt.Sprintf("Reward to '%s'", to)

	txin := newTxInput([]byte{}, -1, data)
	txout := newTxOutput(InitialMiningSubsidy, to)
	tx := newTransaction([]*pb.TXInput{txin}, []*pb.TXOutput{txout})

	return tx
}

func NewTransaction(from, to string, amount int, bc *Blockchain) (*pb.Transaction, error) {
	acc, outputs, err := bc.FindSpendableAmountAndOutputs(from, amount)
	if err != nil {
		return nil, err
	}
	if acc < amount {
		return nil, ErrorNotEnoughBalance
	}

	var (
		txInputs  = make([]*pb.TXInput, 0, len(outputs))
		txOutputs = make([]*pb.TXOutput, 0, 2)
	)

	for id, outs := range outputs {
		txID, err := hex.DecodeString(id)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error of decoding id : %s", id))
		}

		for _, out := range outs {
			in := newTxInput(txID, out, from)
			txInputs = append(txInputs, in)
		}
	}

	txOutputs = append(txOutputs, newTxOutput(amount, to)) // balance change for receiver
	if acc > amount {
		txOutputs = append(txOutputs, newTxOutput(acc-amount, from)) // balance change for sender
	}
	transaction := newTransaction(txInputs, txOutputs)
	return transaction, nil
}
