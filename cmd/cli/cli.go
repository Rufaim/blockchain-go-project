package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Rufaim/blockchain/blockchain"
)

type CLIAppplication struct{}

func NewCLIAppplication() *CLIAppplication {
	return &CLIAppplication{}
}

func (cli *CLIAppplication) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		return
	}

	helpBlockchainFS := flag.NewFlagSet(helpBlockchainCommand, flag.ExitOnError)
	isHelpCall := helpBlockchainFS.Bool(helpBlockchainCommandHFlag, false, "prints this help")

	err := helpBlockchainFS.Parse(os.Args[1:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}

	if *isHelpCall || os.Args[1] == helpBlockchainCommand {
		cli.printUsage()
		return
	}

	switch os.Args[1] {
	case showBlockchainCommand:
		cli.showCommand()
	case balanceOfWalletCommand:
		cli.balanceOfWalletCommand()
	case sendCoinsCommand:
		cli.sendCoinsCommand()
	case createBlockchainCommand:
		cli.createCommand()
	case deleteBlockchainCommand:
		cli.deleteCommand()
	default:
		cli.printUsage()
		return
	}
}

func (cli *CLIAppplication) showCommand() {
	showBlockchainFS := flag.NewFlagSet(showBlockchainCommand, flag.ExitOnError)
	showBlockchainDBAddress := showBlockchainFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := showBlockchainFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsDB(*showBlockchainDBAddress) {
		fmt.Printf("%s is not a valid database\n", *showBlockchainDBAddress)
		return
	}
	if !checkFileExists(*showBlockchainDBAddress) {
		fmt.Printf("Database %s does not exist\n", *showBlockchainDBAddress)
		return
	}
	cli.printChain("", *showBlockchainDBAddress)
}

func (cli *CLIAppplication) balanceOfWalletCommand() {
	balanceOfWalletFS := flag.NewFlagSet(balanceOfWalletCommand, flag.ExitOnError)
	balanceOfWalletFSOF := balanceOfWalletFS.String(balanceOfWalletCommandOFFlag, "", "username to calculate balance")
	balanceOfWalletFSDBAddress := balanceOfWalletFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := balanceOfWalletFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsDB(*balanceOfWalletFSDBAddress) {
		fmt.Printf("%s is not a valid database\n", *balanceOfWalletFSDBAddress)
		return
	}
	if checkFileExists(*balanceOfWalletFSDBAddress) {
		fmt.Printf("Database %s is exist\n", *balanceOfWalletFSDBAddress)
		return
	}
	if len(*balanceOfWalletFSOF) == 0 {
		fmt.Printf("Username is not entered")
		return
	}

	bc, err := blockchain.NewBlockchain(*balanceOfWalletFSDBAddress)
	panicOnError(err)

	txs, err := bc.FindUnspentTransactions(*balanceOfWalletFSOF)
	panicOnError(err)

	balance := 0

	for _, tx := range txs {
		for _, txout := range tx.Outs {
			if blockchain.OutputIsLockedWithKey(txout, *balanceOfWalletFSOF) {
				balance += int(txout.Amount)
			}
		}
	}
	fmt.Printf("Balance for %s is %d\n", *balanceOfWalletFSOF, balance)
}

func (cli *CLIAppplication) sendCoinsCommand() {
	sendCoinsCommandFS := flag.NewFlagSet(sendCoinsCommand, flag.ExitOnError)
	sendCoinsCommandFSFromFlag := sendCoinsCommandFS.String(sendCoinsCommandFromFlag, "", "username for coin sender")
	sendCoinsCommandFSToFlag := sendCoinsCommandFS.String(sendCoinsCommandToFlag, "", "username for coin receiver")
	sendCoinsCommandFSAmountFlag := sendCoinsCommandFS.Int(sendCoinsCommandAmountFlag, -1, "amount of coins")
	sendCoinsCommandFSDBAddress := sendCoinsCommandFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := sendCoinsCommandFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsDB(*sendCoinsCommandFSDBAddress) {
		fmt.Printf("%s is not a valid database\n", *sendCoinsCommandFSDBAddress)
		return
	}
	if checkFileExists(*sendCoinsCommandFSDBAddress) {
		fmt.Printf("Database %s is exist\n", *sendCoinsCommandFSDBAddress)
		return
	}
	if len(*sendCoinsCommandFSFromFlag) == 0 {
		fmt.Printf("Sender username is not entered")
		return
	}
	if len(*sendCoinsCommandFSToFlag) == 0 {
		fmt.Printf("Receiver username is not entered")
		return
	}
	if *sendCoinsCommandFSAmountFlag < 0 {
		fmt.Printf("Receiver username is not entered")
		return
	}

	bc, err := blockchain.NewBlockchain(*sendCoinsCommandFSDBAddress)
	panicOnError(err)

	tx, err := blockchain.NewTransaction(*sendCoinsCommandFSFromFlag, *sendCoinsCommandFSToFlag, *sendCoinsCommandFSAmountFlag, bc)
	_ = tx
	//TODO: send coins implementation
}

func (cli *CLIAppplication) createCommand() {
	createBlockchainFS := flag.NewFlagSet(createBlockchainCommand, flag.ExitOnError)
	createBlockchainFSForceRecreateFlag := createBlockchainFS.Bool(createBlockchainCommandForceRecreateFlag, false, "force recreation of database base")
	createBlockchainDBAddress := createBlockchainFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := createBlockchainFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsDB(*createBlockchainDBAddress) {
		fmt.Printf("%s is not a valid database\n", *createBlockchainDBAddress)
		return
	}
	if checkFileExists(*createBlockchainDBAddress) {
		fmt.Printf("Database %s is exist\n", *createBlockchainDBAddress)
		return
	}

	if checkFileExists(*createBlockchainDBAddress) {
		if *createBlockchainFSForceRecreateFlag {
			if !removeDBFile(*createBlockchainDBAddress) {
				panic("Can not recreate database")
			}
			fmt.Printf("Database %s force recreation!\n", *createBlockchainDBAddress)
		} else {
			fmt.Printf("Database %s is exist\n", *createBlockchainDBAddress)
		}
	}
	bc, err := blockchain.NewBlockchain(*createBlockchainDBAddress)
	panicOnError(err)

	bc.Flush()
}

func (cli *CLIAppplication) deleteCommand() {
	deleteBlockchainFS := flag.NewFlagSet(deleteBlockchainCommand, flag.ExitOnError)
	deleteBlockchainDBAddress := deleteBlockchainFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := deleteBlockchainFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkFileExists(*deleteBlockchainDBAddress) {
		return
	}
	if !checkIsDB(*deleteBlockchainDBAddress) {
		fmt.Printf("%s is not a valid database\n", *deleteBlockchainDBAddress)
		return
	}
	if removeDBFile(*deleteBlockchainDBAddress) {
		fmt.Printf("Database %s removed succesfully!\n", *deleteBlockchainDBAddress)
	}
}
