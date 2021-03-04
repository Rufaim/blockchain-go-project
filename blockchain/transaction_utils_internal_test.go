//+build tests

package blockchain

import (
	"bytes"
	"testing"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
)

func TestInputUsesKey(t *testing.T) {

	w1, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}
	w2, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}

	type args struct {
		input      *pb.TXInput
		pubKeyHash []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Manual conctruction correct", args{&pb.TXInput{
			PubKey: w1.PublicKey,
		}, w1.GetAddressInfo().PublicKeyHash}, true},
		{"Manual construction fail", args{&pb.TXInput{
			PubKey: w2.PublicKey,
		}, w1.GetAddressInfo().PublicKeyHash}, false},
		{"newTxInput correct", args{newTxInput([]byte{}, 0, w2.PublicKey),
			w2.GetAddressInfo().PublicKeyHash}, true},
		{"newTxInput fail", args{newTxInput([]byte{}, 0, w2.PublicKey),
			w1.GetAddressInfo().PublicKeyHash}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InputUsesKey(tt.args.input, tt.args.pubKeyHash); got != tt.want {
				t.Errorf("InputUsesKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputIsLockedWithKey(t *testing.T) {
	w1, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}
	w2, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}

	type args struct {
		output     *pb.TXOutput
		pubKeyHash []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Manual conctruction correct", args{&pb.TXOutput{
			PubKeyHash: w1.GetAddressInfo().PublicKeyHash,
		}, w1.GetAddressInfo().PublicKeyHash}, true},
		{"Manual construction fail", args{&pb.TXOutput{
			PubKeyHash: w2.GetAddressInfo().PublicKeyHash,
		}, w1.GetAddressInfo().PublicKeyHash}, false},
		{"newTxOutput correct", args{newTxOutput(0, w2.GetAddress()),
			w2.GetAddressInfo().PublicKeyHash}, true},
		{"newTxOutput fail", args{newTxOutput(0, w2.GetAddress()),
			w1.GetAddressInfo().PublicKeyHash}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OutputIsLockedWithKey(tt.args.output, tt.args.pubKeyHash); got != tt.want {
				t.Errorf("OutputIsLockedWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTransactionCoinbase(t *testing.T) {
	tests := []struct {
		name string
		tx   *pb.Transaction
		want bool
	}{
		{"Empty TX", &pb.Transaction{}, false},
		{"Manual conctruction correct", &pb.Transaction{
			Inps: []*pb.TXInput{{
				OutId: int32(-1),
			}},
			Outs: []*pb.TXOutput{},
		}, true},
		{"Manual conctruction fail 1", &pb.Transaction{
			Inps: []*pb.TXInput{},
			Outs: []*pb.TXOutput{},
		}, false},
		{"Manual conctruction fail 2", &pb.Transaction{
			Inps: []*pb.TXInput{{
				Id:    []byte("165165151"),
				OutId: int32(2),
			}},
			Outs: []*pb.TXOutput{},
		}, false},
		{"Manual conctruction fail 3", &pb.Transaction{
			Inps: []*pb.TXInput{{
				Id:    []byte("165165151"),
				OutId: int32(-1),
			}},
			Outs: []*pb.TXOutput{},
		}, false},
		{"NewCoinbaseTX constructor", NewCoinbaseTX([]byte("123456789")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTransactionCoinbase(tt.tx); got != tt.want {
				t.Errorf("OutputIsLockedWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllTransactionInputsIds(t *testing.T) {
	w, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Wallet generation error")
	}
	tests := []struct {
		name string
		tx   *pb.Transaction
		want [][]byte
	}{
		{"Empty TX", &pb.Transaction{}, [][]byte{}},
		{"Manual construction 1", &pb.Transaction{
			Inps: []*pb.TXInput{{Id: []byte("123")}},
		}, [][]byte{[]byte("123")}},
		{"Manual construction 2", &pb.Transaction{
			Inps: []*pb.TXInput{{Id: []byte("123")}, {Id: []byte("234")}},
		}, [][]byte{[]byte("123"), []byte("234")}},
		{"Coinbase", NewCoinbaseTX(w.GetAddress()), [][]byte{[]byte("")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAllTransactionInputsIds(tt.tx)
			checkedWant := make([]bool, len(tt.want))
			checkedGot := make([]bool, len(got))

			for i, gotid := range got {
				for j, wantid := range tt.want {
					if bytes.Compare(gotid, wantid) == 0 {
						if checkedWant[j] {
							t.Errorf("Id returned multiple times")
						}
						checkedGot[i] = true
						checkedWant[j] = true
						break
					}
				}
			}
			for _, c := range checkedGot {
				if !c {
					t.Errorf("Additional ids returned")
				}
			}
			for _, c := range checkedWant {
				if !c {
					t.Errorf("Some ids are missed")
				}
			}
		})
	}
}
