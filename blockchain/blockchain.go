package blockchain

import (
	"bytes"
	"encoding/hex"
	"time"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

//Blockchain is a primary structure containig a chain
type Blockchain struct {
	currentTop *block
	db         *bolt.DB
}

//MineBlock is used to pass a new data to blockchain
//anything implementing BlockData interface is a valid data
func (bc *Blockchain) MineBlock(transactions []*pb.Transaction) ([]byte, error) {
	if len(bc.currentTop.Hash) == 0 {
		bc.currentTop.setHash()
	}
	newBlock := newBlock(transactions, bc.currentTop.Hash)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		repr, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		bucket.Put(newBlock.Hash, repr)
		bucket.Put([]byte(keyTopBlockHash), newBlock.Hash)
		return nil
	})
	if err != nil {
		return []byte{}, err
	}

	bc.currentTop = newBlock
	return newBlock.Hash, nil
}

//NewBlockchain is a blockchain constructor
func NewBlockchain(dbpath string, founderAddress []byte) (*Blockchain, error) {
	db, err := bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, errors.Wrap(err, "Database opening failure")
	}

	var top *block
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		if bucket != nil {
			topHash := bucket.Get([]byte(keyTopBlockHash))
			top, err = DeserializeBlock(bucket.Get(topHash))
			if err != nil {
				return err
			}
			return nil
		}

		bucket, err := tx.CreateBucket([]byte(blocksBucketName))
		cbtx := NewCoinbaseTX(founderAddress)
		top := newGenesisBlock(cbtx)
		repr, err := top.Serialize()
		if err != nil {
			return err
		}
		bucket.Put(top.Hash, repr)
		bucket.Put([]byte(keyTopBlockHash), top.Hash)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		currentTop: top,
		db:         db,
	}, nil
}

func (bc *Blockchain) Flush() error {
	return bc.db.Sync()
}

func (bc *Blockchain) Finalize() error {
	return bc.db.Close()
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.currentTop.Hash,
		db:          bc.db,
	}
}

//FindTransactionsByID search for a set of transactions in the Blockchain
// returns ErrorTransactionsNotFound is even one transaction id does not exist
func (bc *Blockchain) FindTransactionsByID(IDs [][]byte) ([]*pb.Transaction, error) {
	bci := bc.Iterator()
	type boolID struct {
		isFound bool
		id      []byte
	}
	boolMapper := make([]*boolID, 0, len(IDs))
	for _, id := range IDs {
		boolMapper = append(boolMapper, &boolID{false, id})
	}
	result := make([]*pb.Transaction, 0, len(IDs))

	for {

		block, err := bci.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			for _, mid := range boolMapper {
				if !mid.isFound {
					if bytes.Compare(tx.Id, mid.id) == 0 {
						result = append(result, tx)
						mid.isFound = true
					}
				}
			}
		}

		if block.IsGenesis() {
			break
		}

		stopKey := true
		for _, mid := range boolMapper {
			if !mid.isFound {
				stopKey = false
				break
			}
		}
		if stopKey {
			return result, nil
		}
	}
	for _, mid := range boolMapper {
		if !mid.isFound {
			return nil, ErrorTransactionsNotFound
		}
	}
	return result, nil
}

//FindUnspentTransactions returns a set of transaction that have not been closed.
func (bc *Blockchain) FindUnspentTransactions(address []byte) ([]*pb.Transaction, error) {
	var unspentTXs []*pb.Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	wi := wallet.GetAddressInfo(address)

	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.Id)

		Outputs:
			for outIdx, out := range tx.Outs {
				// If output is not spent ...
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// and if it is target output ...
				if OutputIsLockedWithKey(out, wi.PublicKeyHash) {
					// we take it into account
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !isTransactionCoinbase(tx) {
				for _, in := range tx.Inps {
					if InputUsesKey(in, wi.PublicKeyHash) {
						// remember, id of the input is an id of transaction whose output it closes
						inTxID := hex.EncodeToString(in.Id)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], int(in.OutId))
					}
				}
			}
		}

		if block.IsGenesis() {
			break
		}
	}

	return unspentTXs, nil
}

//FindSpendableAmountAndOutputs returns a set of unfinished transaction outputs
// with a sum more than provided amount. if amount of less then zero provided, than
// total sum of all unspent transactions returened
func (bc *Blockchain) FindSpendableAmountAndOutputs(address []byte, amount int) (int, map[string][]int, error) {
	unspentOutputs := make(map[string][]int)
	transactionBalance := 0
	unspentTransactions, err := bc.FindUnspentTransactions(address)
	wi := wallet.GetAddressInfo(address)

	if err != nil {
		return transactionBalance, unspentOutputs, err
	}

	for _, tx := range unspentTransactions {
		txID := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.Outs {
			if transactionBalance >= amount && amount >= 0 {
				return transactionBalance, unspentOutputs, nil
			}
			if OutputIsLockedWithKey(out, wi.PublicKeyHash) {
				transactionBalance += int(out.Amount)
				if amount > 0 {
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
	}

	return transactionBalance, unspentOutputs, nil
}
