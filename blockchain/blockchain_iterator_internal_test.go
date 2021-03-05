package blockchain

import (
	"bytes"
	"encoding/hex"
	"testing"

	pb "github.com/Rufaim/blockchain/message"
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
			}
			for key, val := range tt.want {
				tx, ok := got[key]
				if !ok {
					t.Errorf("BlockchainIterator.FindTransactionsByID() not returned transaction of id %s", string(val.Id))
				}
				if bytes.Compare(tx.Id, val.Id) != 0 {
					t.Errorf("BlockchainIterator.FindTransactionsByID() returned transaction for id %s", string(val.Id))
				}
			}
		})
	}
}
