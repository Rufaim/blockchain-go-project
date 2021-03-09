//+build tests

package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
)

type mockBlockchainIterator struct {
	blockList []*block
	curretPos int
}

func (mb *mockBlockchainIterator) Next() (*block, error) {
	mb.curretPos++
	return mb.blockList[mb.curretPos], nil
}

func TestBlockchainIterator_FindTransactionsByID(t *testing.T) {
	blockList := []*block{
		{pb.Block{
			PrevHash: []byte("0"), // we know taht genesis does not have previous hash at all
			Transactions: []*pb.Transaction{
				{Id: []byte("123")},
				{Id: []byte("234")},
			},
		}},
		{pb.Block{
			PrevHash: []byte("0"),
			Transactions: []*pb.Transaction{
				{Id: []byte("345")},
			},
		}},
		{pb.Block{
			PrevHash: []byte("0"),
			Transactions: []*pb.Transaction{
				{Id: []byte("456")},
				{Id: []byte("567")},
				{Id: []byte("678")},
			},
		}},
		{},
	}

	tests := []struct {
		name    string
		IDs     [][]byte
		want    map[string]*pb.Transaction
		wantErr bool
	}{
		{"One transaction", [][]byte{[]byte("123")}, map[string]*pb.Transaction{hex.EncodeToString([]byte("123")): blockList[0].Transactions[0]}, false},
		{"Invalid request", [][]byte{[]byte("aaa")}, nil, true},
		{"Valid plus invalid", [][]byte{[]byte("345"), []byte("aaa"), []byte("567")}, nil, true},
		{"Multiple transactions", [][]byte{[]byte("678"), []byte("345"), []byte("123")},
			map[string]*pb.Transaction{hex.EncodeToString([]byte("123")): blockList[0].Transactions[0],
				hex.EncodeToString([]byte("678")): blockList[2].Transactions[2],
				hex.EncodeToString([]byte("345")): blockList[1].Transactions[0]},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &mockBlockchainIterator{
				blockList: blockList,
				curretPos: -1,
			}
			got, err := FindTransactionsByID(bi, tt.IDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("BlockchainIterator.FindTransactionsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("BlockchainIterator.FindTransactionsByID() returned wrong number of results")
				return
			}
			for key, val := range tt.want {
				tx, ok := got[key]
				if !ok {
					t.Errorf("BlockchainIterator.FindTransactionsByID() not returned transaction of id %s", string(val.Id))
					return
				}
				if bytes.Compare(tx.Id, val.Id) != 0 {
					t.Errorf("BlockchainIterator.FindTransactionsByID() returned transaction for id %s", string(val.Id))
					return
				}
			}
		})
	}
}

func TestBlockchainIterator_FindUnspentTransactions(t *testing.T) {
	pubKey1 := []byte("aaa")
	pubKey1Hash := wallet.HashPubKey(pubKey1)
	pubKey2 := []byte("bbb")
	pubKey2Hash := wallet.HashPubKey(pubKey2)
	blockList := []*block{
		{pb.Block{
			PrevHash: []byte("0"), // we know taht genesis does not have previous hash at all
			Transactions: []*pb.Transaction{
				{Id: []byte("123"),
					Inps: []*pb.TXInput{{Id: []byte("345"), OutId: 0, PubKey: pubKey1}},
					Outs: []*pb.TXOutput{{Amount: 2, PubKeyHash: []byte("jnea")}, {Amount: 8, PubKeyHash: pubKey1Hash}},
				},
				// {Id: []byte("234")},
			},
		}},
		{pb.Block{
			PrevHash: []byte("0"),
			Transactions: []*pb.Transaction{
				{Id: []byte("456"),
					Inps: []*pb.TXInput{{Id: []byte("456"), OutId: 0, PubKey: []byte("jnea")}},
					Outs: []*pb.TXOutput{{Amount: 2, PubKeyHash: pubKey2Hash}},
				},
				{Id: []byte("345"), // coinbase key 1
					Inps: []*pb.TXInput{{OutId: -1}},
					Outs: []*pb.TXOutput{{Amount: 10, PubKeyHash: pubKey1Hash}},
				},
			},
		}},
		{pb.Block{
			PrevHash: []byte("0"),
			Transactions: []*pb.Transaction{
				{Id: []byte("456"),
					Inps: []*pb.TXInput{{Id: []byte("678"), OutId: 2, PubKey: []byte("afdbmad")}},
					Outs: []*pb.TXOutput{{Amount: 3, PubKeyHash: []byte("jnea")}},
				},
				{Id: []byte("345"), // coinbase key 1
					Inps: []*pb.TXInput{{OutId: -1}},
					Outs: []*pb.TXOutput{{Amount: 10, PubKeyHash: pubKey1Hash}},
				},
				{Id: []byte("567"), // coinbase key 2
					Inps: []*pb.TXInput{{OutId: -1}},
					Outs: []*pb.TXOutput{{Amount: 10, PubKeyHash: pubKey2Hash}},
				},
			},
		}},
		{},
	}

	tests := []struct {
		name string
		key  []byte
		want [][]byte
	}{
		{"key 1", pubKey1Hash, [][]byte{[]byte("123"), []byte("345")}},
		{"key 2", pubKey2Hash, [][]byte{[]byte("456"), []byte("567")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &mockBlockchainIterator{
				blockList: blockList,
				curretPos: -1,
			}
			got, err := FindUnspentTransactions(bi, tt.key)

			for _, v := range got {
				fmt.Println(string(v.Id))
			}
			fmt.Println()
			for _, v := range tt.want {
				fmt.Println(string(v))
			}

			if err != nil {
				t.Errorf("FindUnspentTransactions() error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("FindUnspentTransactions() returned wrong number of results")
				return
			}
			for i := range tt.want {
				if bytes.Compare(tt.want[i], got[i].Id) != 0 {
					t.Errorf("FindUnspentTransactions() returned transaction for id %s", string(got[i].Id))
					return
				}
			}
		})
	}
}
