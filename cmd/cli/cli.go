package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Rufaim/blockchain/blockchain"
	pb "github.com/Rufaim/blockchain/message"
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
	case addDataCommand:
		cli.addDataCommand()
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

func (cli *CLIAppplication) addDataCommand() {
	addDataFS := flag.NewFlagSet(addDataCommand, flag.ExitOnError)
	addDataFSData := addDataFS.String(addDataCommandDataFlag, "", "data to add into a blockchain")
	addDataFSDBAddress := addDataFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := addDataFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if len(*addDataFSData) == 0 {
		fmt.Printf("Data is not entered")
		return
	}

	bc, err := blockchain.NewBlockchain(*addDataFSDBAddress)
	panicOnError(err)

	tx := blockchain.NewTransaction([]byte(*addDataFSData))
	hash, err := bc.MineBlock([]*pb.Transaction{tx})
	panicOnError(err)
	fmt.Printf("Data added!\nHASH: %x\n", hash)
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
