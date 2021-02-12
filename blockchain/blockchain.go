package blockchain

import (
	"encoding/hex"
	"time"

	pb "github.com/Rufaim/blockchain/message"
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
func NewBlockchain(dbpath string) (*Blockchain, error) {
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
		cbtx := NewCoinbaseTX(genesisCoinbaseData)
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

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.currentTop.Hash,
		db:          bc.db,
	}
}

//FindUnspentTransactions returns a set of transaction that have not been closed.
func (bc *Blockchain) FindUnspentTransactions(pubKey string) ([]*pb.Transaction, error) {
	var unspentTXs []*pb.Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

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
				if outputIsLockedWithKey(out, pubKey) {
					// we take it into account
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !isTransactionCoinbase(tx) {
				for _, in := range tx.Inps {
					if inputUsesKey(in, pubKey) {
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
func (bc *Blockchain) FindSpendableAmountAndOutputs(pubKey string, amount int) (int, map[string][]int, error) {
	unspentOutputs := make(map[string][]int)
	transactionBalance := 0
	unspentTransactions, err := bc.FindUnspentTransactions(pubKey)

	if err != nil {
		return transactionBalance, unspentOutputs, err
	}

	for _, tx := range unspentTransactions {
		txID := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.Outs {
			if transactionBalance >= amount && amount >= 0 {
				return transactionBalance, unspentOutputs, nil
			}
			if outputIsLockedWithKey(out, pubKey) {
				transactionBalance += int(out.Amount)
				if amount > 0 {
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
	}

	return transactionBalance, unspentOutputs, nil
}
