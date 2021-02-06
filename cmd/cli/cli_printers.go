package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Rufaim/blockchain/blockchain"
)

func (*CLIAppplication) printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s %s -%s [-%s dbpath] -- adds data to a chain\n", os.Args[0], addDataCommand, addDataCommandDataFlag, DBAddressFlag)
	fmt.Printf("%s %s [-%s dbpath] -- prints the blockchain from the last to genesis\n", os.Args[0],
		showBlockchainCommand, DBAddressFlag)
	fmt.Printf("%s %s [-%s] [-%s dbpath] -- creates new blockhain\n", os.Args[0],
		createBlockchainCommand, createBlockchainCommandForceRecreateFlag, DBAddressFlag)
	fmt.Printf("%s %s [-%s dbpath] -- deletes given blockchain\n", os.Args[0],
		deleteBlockchainCommand, DBAddressFlag)
	fmt.Printf("%s [-%s | %s] -- prints this help\n", os.Args[0], helpBlockchainCommandHFlag, helpBlockchainCommand)
}

func (*CLIAppplication) printChain(address, dbPath string) {
	bc, err := blockchain.NewBlockchain(dbPath)
	if err != nil {
		panic(err)
	}
	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Prev: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Printf("Data: %s\n", string(block.Transactions[0].Data))
		fmt.Println()

		if block.IsGenesis() {
			break
		}
	}
}
