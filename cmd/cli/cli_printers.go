package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Rufaim/blockchain/blockchain"
)

func (*CLIAppplication) printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s %s -%s username [-%s dbpath] -- prints a balance for a given user\n", os.Args[0], balanceOfWalletCommand, balanceOfWalletCommandOFFlag, DBAddressFlag)
	fmt.Printf("%s %s -%s sender -%s receiver -%s amount [-%s dbpath] -- prints a balance for a given user\n", os.Args[0], sendCoinsCommand, sendCoinsCommandFromFlag, sendCoinsCommandToFlag, sendCoinsCommandAmountFlag, DBAddressFlag)
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
	panicOnError(err)
	defer bc.Finalize()

	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Prev: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Println("Transactions:")
		for _, tx := range block.Transactions {
			fmt.Println("\tInputs:")
			for _, in := range tx.Inps {
				fmt.Printf("\t\tOutID: %d; PubKey: %s\n", in.OutId, in.PubKey)
			}
			fmt.Println("\tOutputs:")
			for _, out := range tx.Outs {
				fmt.Printf("\t\tAmount: %d; PubKey: %s\n", out.Amount, out.PubKey)
			}
		}

		fmt.Println()

		if block.IsGenesis() {
			break
		}
	}
}
