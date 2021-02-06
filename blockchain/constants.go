package blockchain

import (
	"crypto/sha256"
	"math/big"
)

//// Proof of Work constants
const hashTargetBits = 24 //remember <target> = <desired amount of zeros in hash> * 4

var hashTargetValue *big.Int

func init() {
	hashTargetValue = big.NewInt(1)
	hashTargetValue.Lsh(hashTargetValue, uint(sha256.Size*8-hashTargetBits))
}

//// Database constants
const (
	blocksBucketName    = "blocks"
	keyTopBlockHash     = "top"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

const miningSubsidy = 10
