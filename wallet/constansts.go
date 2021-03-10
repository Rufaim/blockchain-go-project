package wallet

import "errors"

const (
	Version            = byte(0x00)
	PublicKeyHashSize  = 20
	AddressChecksumLen = 4
)

var ErrorWalletDoesNotExist = errors.New("wallet does not exist")
