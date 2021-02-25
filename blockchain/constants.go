package blockchain

import (
	"crypto/sha256"
	"errors"
	"math/big"
)

//// Proof of Work constants
const hashTargetBits = 16 //remember <target> = <desired amount of zeros in hash> * 4

var hashTargetValue *big.Int

func init() {
	hashTargetValue = big.NewInt(1)
	hashTargetValue.Lsh(hashTargetValue, uint(sha256.Size*8-hashTargetBits))
}

//// Database constants
const (
	blocksBucketName    = "blocks"
	keyTopBlockHash     = "top"
	genesisCoinbaseData = "The Times 27/Jul/2020 Chancellor on brink of second bailout for banks"
)

const InitialMiningSubsidy = 10

var (
	ErrorNotEnoughBalance      = errors.New("Not enougth balance for transaction")
	ErrorTransactionsNotFound  = errors.New("Transactions are not found")
	ErrorInvalidTransaction    = errors.New("Transactions is invalid")
	ErrorTransactionCopyFailed = errors.New("Transactions if failed to copy")
)
