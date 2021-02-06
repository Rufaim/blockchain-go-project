package blockchain

import (
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
		cbtx := NewTransaction([]byte(genesisCoinbaseData))
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
