package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
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

func newTransaction(inps []*pb.TXInput, outs []*pb.TXOutput) *pb.Transaction {
	tx := &pb.Transaction{
		Inps: inps,
		Outs: outs,
	}

	tx.Id = hashTransaction(tx)
	return tx
}

func SignTransactionWithWallet(tx *pb.Transaction, w *wallet.Wallet, refTxs map[string]*pb.Transaction) error {
	if isTransactionCoinbase(tx) {
		return nil
	}

	txCopy, ok := proto.Clone(tx).(*pb.Transaction)
	if !ok {
		return ErrorTransactionCopyFailed
	}

	for i := range tx.Inps {
		txCopy.Inps[i].Signature = nil
		txCopy.Inps[i].PubKey = nil
	}

	for i, inp := range tx.Inps {
		txCopy.Inps[i].PubKey = refTxs[hex.EncodeToString(inp.Id)].Outs[int(inp.OutId)].PubKeyHash
		txCopy.Id = hashTransaction(txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &w.PrivateKey, txCopy.Id)
		if err != nil {
			return err
		}

		tx.Inps[i].Signature = append(r.Bytes(), s.Bytes()...)
	}

	return nil
}

func VerifyTransaction(tx *pb.Transaction, refTxs map[string]*pb.Transaction) (bool, error) {
	if isTransactionCoinbase(tx) {
		return true, nil
	}

	txCopy, ok := proto.Clone(tx).(*pb.Transaction)
	if !ok {
		return false, ErrorTransactionCopyFailed
	}
	for i := range tx.Inps {
		txCopy.Inps[i].Signature = nil
		txCopy.Inps[i].PubKey = nil
	}

	curve := elliptic.P256()

	for i, inp := range tx.Inps {
		lenHalf := len(inp.PubKey) / 2
		x := new(big.Int).SetBytes(inp.PubKey[:lenHalf])
		y := new(big.Int).SetBytes(inp.PubKey[lenHalf:])
		recPubKey := ecdsa.PublicKey{curve, x, y}

		lenHalf = len(inp.Signature) / 2
		r := new(big.Int).SetBytes(inp.Signature[:lenHalf])
		s := new(big.Int).SetBytes(inp.Signature[lenHalf:])

		txCopy.Inps[i].Signature = nil
		txCopy.Inps[i].PubKey = refTxs[hex.EncodeToString(inp.Id)].Outs[int(inp.OutId)].PubKeyHash
		txCopy.Id = hashTransaction(txCopy)

		if !ecdsa.Verify(&recPubKey, txCopy.Id, r, s) {
			return false, nil
		}
	}

	return true, nil
}

func NewCoinbaseTX(to []byte) *pb.Transaction {
	txin := newTxInput([]byte{}, -1, []byte(genesisCoinbaseData), []byte{})
	txout := newTxOutput(InitialMiningSubsidy, to)
	tx := newTransaction([]*pb.TXInput{txin}, []*pb.TXOutput{txout})

	return tx
}

func NewTransaction(from, to []byte, amount int, bc *Blockchain, ws *wallet.WalletSet) (*pb.Transaction, error) {
	acc, outputs, err := bc.FindSpendableAmountAndOutputs(from, amount)
	if err != nil {
		return nil, err
	}
	if acc < amount {
		return nil, ErrorNotEnoughBalance
	}

	wfrom, err := ws.GetWalletByAddress(string(from))
	if err != nil {
		return nil, err
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
			in := newTxInput(txID, out, wfrom.PublicKey, []byte{})
			txInputs = append(txInputs, in)
		}
	}

	txOutputs = append(txOutputs, newTxOutput(amount, to)) // balance change for receiver
	if acc > amount {
		txOutputs = append(txOutputs, newTxOutput(acc-amount, from)) // balance change for sender
	}
	transaction := newTransaction(txInputs, txOutputs)
	err = bc.SignTransaction(transaction, wfrom)
	return transaction, err
}
