package blockchain

import (
	"bytes"
	"encoding/hex"

	pb "github.com/Rufaim/blockchain/message"
	"github.com/boltdb/bolt"
)

//Unified interface for blockchain iterators
type BlockchainIterable interface {
	Next() (*block, error)
}

//BlockchainIterator is a top to genesis iterator over a blockchain
//Usage:
//bci := <iterator construction>
//for {
//		block, err := bci.Next()
//		if err != nil {
//			return nil, err
//		}
//		...
//}
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

//FindTransactionsByID search for a set of transactions in the Blockchain
// returns ErrorTransactionsNotFound is even one transaction id does not exist
func FindTransactionsByID(bi BlockchainIterable, IDs [][]byte) (map[string]*pb.Transaction, error) {
	type boolID struct {
		isFound bool
		id      []byte
	}
	boolMapper := make([]*boolID, 0, len(IDs))
	for _, id := range IDs {
		boolMapper = append(boolMapper, &boolID{false, id})
	}
	result := make(map[string]*pb.Transaction, len(IDs))

	for {
		block, err := bi.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			for _, mid := range boolMapper {
				if !mid.isFound {
					if bytes.Compare(tx.Id, mid.id) == 0 {
						result[hex.EncodeToString(tx.Id)] = tx
						mid.isFound = true
					}
				}
			}
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

		if block.IsGenesis() {
			break
		}
	}
	return nil, ErrorTransactionsNotFound
}

//FindUnspentTransactions returns a set of transaction that have not been closed.
func FindUnspentTransactions(bi BlockchainIterable, pubKeyHash []byte) ([]*pb.Transaction, error) {
	type spent struct {
		idx       int
		accounted bool
	}

	var unspentTXs []*pb.Transaction
	spentTXOs := make(map[string][]*spent)

	for {
		block, err := bi.Next()
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
						if !spentOut.accounted && spentOut.idx == outIdx {
							spentOut.accounted = true
							continue Outputs
						}
					}
				}
				// and if it is target output ...
				if OutputIsLockedWithKey(out, pubKeyHash) {
					// we take it into account
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !isTransactionCoinbase(tx) {
				for _, in := range tx.Inps {
					if InputUsesKey(in, pubKeyHash) {
						// remember, id of the input is an id of transaction whose output it closes
						inTxID := hex.EncodeToString(in.Id)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], &spent{idx: int(in.OutId)})
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
