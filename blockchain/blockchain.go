package blockchain

import (
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

	for _, tx := range transactions {
		vres, err := bc.VerifyTransaction(tx)
		if err != nil {
			panic(err)
		}
		if !vres {
			return nil, ErrorInvalidTransaction
		}
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

func (bc *Blockchain) Iterator() BlockchainIterable {
	return &BlockchainIterator{
		currentHash: bc.currentTop.Hash,
		db:          bc.db,
	}
}

func (bc *Blockchain) SignTransaction(tx *pb.Transaction, w *wallet.Wallet) error {
	if isTransactionCoinbase(tx) {
		return nil
	}
	refTxIds := getAllTransactionInputsIds(tx)
	refTXs, err := FindTransactionsByID(bc.Iterator(), refTxIds)

	if err != nil {
		return err
	}

	return SignTransactionWithWallet(tx, w, refTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *pb.Transaction) (bool, error) {
	if isTransactionCoinbase(tx) {
		return true, nil
	}
	refTxIds := getAllTransactionInputsIds(tx)
	refTXs, err := FindTransactionsByID(bc.Iterator(), refTxIds)
	if err != nil {
		return false, err
	}

	return VerifyTransaction(tx, refTXs)
}

//FindUnspentTransactions returns a set of transaction that have not been closed.
func (bc *Blockchain) FindUnspentTransactions(address []byte) ([]*pb.Transaction, error) {
	wi := wallet.GetAddressInfo(address)
	return FindUnspentTransactions(bc.Iterator(), wi.PublicKeyHash)
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
