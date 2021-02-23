package cli

import (
	"os"
	"strings"

	"github.com/Rufaim/blockchain/wallet"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func removeDBFile(path string) bool {
	if !checkIsExt(path, ".db") {
		return false
	}

	panicOnError(os.Remove(path))
	return !checkFileExists(path)
}

func checkIsExt(path, suffix string) bool {
	if strings.HasSuffix(path, suffix) {
		return true
	}
	return false
}

func checkWalletAddress(addr string) bool {
	return wallet.IsValidAddress([]byte(addr))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
