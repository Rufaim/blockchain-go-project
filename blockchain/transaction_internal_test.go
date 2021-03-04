//+build tests

package blockchain

import (
	"encoding/hex"
	"testing"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
)

func TestSignVerifyTransaction(t *testing.T) {
	w1, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}
	w2, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}

	cbtx := NewCoinbaseTX(w1.GetAddress())
	tx_ := newTransaction([]*pb.TXInput{newTxInput([]byte("123"), 2, w2.PublicKey)},
		[]*pb.TXOutput{newTxOutput(4, w1.GetAddress())})

	tx := newTransaction([]*pb.TXInput{newTxInput(tx_.Id, 0, w1.PublicKey), newTxInput(cbtx.Id, 0, w1.PublicKey)},
		[]*pb.TXOutput{newTxOutput(5, w2.GetAddress())})
	refTxIDs := make(map[string]*pb.Transaction)
	refTxIDs[hex.EncodeToString(cbtx.Id)] = cbtx
	refTxIDs[hex.EncodeToString(tx_.Id)] = tx_

	err = SignTransactionWithWallet(tx, w1, refTxIDs)
	if err != nil {
		t.Fatalf("Signing Transaction error: %s", err.Error())
	}

	verified, err := VerifyTransaction(tx, refTxIDs)
	if err != nil {
		t.Fatalf("Verifying Transaction error: %s", err.Error())
	}
	if !verified {
		t.Errorf("Signed artificial transaction expected to be correct")
	}
}
