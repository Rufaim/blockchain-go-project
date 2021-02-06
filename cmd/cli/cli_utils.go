package cli

import (
	"os"
	"strings"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func removeDBFile(path string) bool {
	if !checkIsDB(path) {
		return false
	}

	panicOnError(os.Remove(path))
	return !checkFileExists(path)
}

func checkIsDB(path string) bool {
	if strings.HasSuffix(path, ".db") {
		return true
	}
	return false
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
