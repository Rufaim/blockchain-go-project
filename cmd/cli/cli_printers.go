package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Rufaim/blockchain/base58"
	"github.com/Rufaim/blockchain/blockchain"
	"github.com/Rufaim/blockchain/wallet"
)

func (*CLIAppplication) printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s %s [-%s] [-%s wspath] -- shows all adresses or create a new one\n", os.Args[0],
		walletOperationsCommand, walletOperationsCommandNewFlag, WSAddressFlag)
	fmt.Printf("%s %s -%s address [-%s dbpath] -- prints a balance for a given user address\n", os.Args[0],
		balanceOfWalletCommand, balanceOfWalletCommandOFFlag, DBAddressFlag)
	fmt.Printf("%s %s -%s sender -%s receiver -%s amount [-%s wspath] [-%s dbpath] -- transfer coins from sender to receiver\n", os.Args[0],
		sendCoinsCommand, sendCoinsCommandFromFlag, sendCoinsCommandToFlag, sendCoinsCommandAmountFlag, WSAddressFlag, DBAddressFlag)
	fmt.Printf("%s %s [-%s dbpath] -- prints the blockchain from the last to genesis\n", os.Args[0],
		showBlockchainCommand, DBAddressFlag)
	fmt.Printf("%s %s -%s [-%s] [-%s dbpath] -- creates new blockhain\n", os.Args[0],
		createBlockchainCommand, createBlockchainCommandAddressFlag, createBlockchainCommandForceRecreateFlag, DBAddressFlag)
	fmt.Printf("%s %s [-%s dbpath] -- deletes given blockchain\n", os.Args[0],
		deleteBlockchainCommand, DBAddressFlag)
	fmt.Printf("%s [-%s | %s] -- prints this help\n", os.Args[0], helpBlockchainCommandHFlag, helpBlockchainCommand)
}

func (*CLIAppplication) printChain(address, dbPath string) {
	bc, err := blockchain.NewBlockchain(dbPath, []byte{})
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
				hash := string(base58.Base58Encode(wallet.HashPubKey(in.PubKey)))
				fmt.Printf("\t\tOutID: %d; PubKeyHash: %s\n", in.OutId, hash)
			}
			fmt.Println("\tOutputs:")
			for _, out := range tx.Outs {
				hash := string(base58.Base58Encode(out.PubKeyHash))
				fmt.Printf("\t\tAmount: %d; PubKeyHash: %s\n", out.Amount, hash)
			}
		}

		fmt.Println()

		if block.IsGenesis() {
			break
		}
	}
}
