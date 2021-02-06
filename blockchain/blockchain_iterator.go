package blockchain

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next return a block that prevhash field of the current block is points to
func (bi *BlockchainIterator) Next() (*block, error) {
	var block *block

	err := bi.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		var err error
		block, err = DeserializeBlock(bucket.Get([]byte(bi.currentHash)))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if !block.IsGenesis() {
		bi.currentHash = block.PrevHash
	}

	return block, nil
}
