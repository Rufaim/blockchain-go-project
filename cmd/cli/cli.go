package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Rufaim/blockchain/blockchain"
	pb "github.com/Rufaim/blockchain/message"
	"github.com/Rufaim/blockchain/wallet"
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
	case walletOperationsCommand:
		cli.walletOperationsCommand()
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
	if !checkIsExt(*showBlockchainDBAddress, ".db") {
		fmt.Printf("%s is not a valid database\n", *showBlockchainDBAddress)
		return
	}
	if !checkFileExists(*showBlockchainDBAddress) {
		fmt.Printf("Database %s does not exist\n", *showBlockchainDBAddress)
		return
	}
	cli.printChain("", *showBlockchainDBAddress)
}

func (cli *CLIAppplication) walletOperationsCommand() {
	walletOperationsFS := flag.NewFlagSet(balanceOfWalletCommand, flag.ExitOnError)
	walletOperationsFSNew := walletOperationsFS.Bool(walletOperationsCommandNewFlag, false, "flag to create new wallet and print its value")
	walletOperationsFSWSAddress := walletOperationsFS.String(WSAddressFlag, defaultWalletSetPath, "provides path to a specific wallet_set.ws")

	err := walletOperationsFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsExt(*walletOperationsFSWSAddress, ".ws") {
		fmt.Printf("%s is not a valid wallet set\n", *walletOperationsFSWSAddress)
		return
	}

	ws := wallet.NewWalletSet()
	if checkFileExists(*walletOperationsFSWSAddress) {
		ws.LoadFromFile(*walletOperationsFSWSAddress)
		panicOnError(err)
	} else if !*walletOperationsFSNew {
		fmt.Printf("Wallet set %s does not exist\n", *walletOperationsFSWSAddress)
		return
	}

	if *walletOperationsFSNew {
		address, err := ws.CreateWallet()
		panicOnError(err)
		fmt.Printf("New address: %s\n", address)
		ws.SaveToFile(*walletOperationsFSWSAddress)
		return
	}

	for _, address := range ws.GetAllAddresses() {
		fmt.Println(address)
	}
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
	if !checkIsExt(*balanceOfWalletFSDBAddress, ".db") {
		fmt.Printf("%s is not a valid database\n", *balanceOfWalletFSDBAddress)
		return
	}
	if !checkFileExists(*balanceOfWalletFSDBAddress) {
		fmt.Printf("Database %s does not exist\n", *balanceOfWalletFSDBAddress)
		return
	}
	if len(*balanceOfWalletFSOF) == 0 {
		fmt.Printf("Username is not entered")
		return
	}

	bc, err := blockchain.NewBlockchain(*balanceOfWalletFSDBAddress, []byte{})
	panicOnError(err)
	defer bc.Finalize()

	txs, err := bc.FindUnspentTransactions([]byte(*balanceOfWalletFSOF))
	panicOnError(err)

	balance := 0
	wi := wallet.GetAddressInfo([]byte(*balanceOfWalletFSOF))
	for _, tx := range txs {
		for _, txout := range tx.Outs {
			if blockchain.OutputIsLockedWithKey(txout, wi.PublicKeyHash) {
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
	sendCoinsCommandFSWSAddress := sendCoinsCommandFS.String(WSAddressFlag, defaultWalletSetPath, "provides path to a specific wallet set file")
	sendCoinsCommandFSDBAddress := sendCoinsCommandFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := sendCoinsCommandFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsExt(*sendCoinsCommandFSDBAddress, ".db") {
		fmt.Printf("%s is not a valid database\n", *sendCoinsCommandFSDBAddress)
		return
	}
	if !checkFileExists(*sendCoinsCommandFSDBAddress) {
		fmt.Printf("Database %s does not exist\n", *sendCoinsCommandFSDBAddress)
		return
	}
	if !checkIsExt(*sendCoinsCommandFSWSAddress, ".ws") {
		fmt.Printf("%s is not a valid wallet set\n", *sendCoinsCommandFSWSAddress)
		return
	}
	if !checkFileExists(*sendCoinsCommandFSWSAddress) {
		fmt.Printf("Wallet set %s does not exist\n", *sendCoinsCommandFSWSAddress)
		return
	}

	ws := wallet.NewWalletSet()
	panicOnError(ws.LoadFromFile(*sendCoinsCommandFSWSAddress))

	if len(*sendCoinsCommandFSFromFlag) == 0 {
		fmt.Printf("Sender address is not entered")
		return
	}
	if !checkWalletAddress(*sendCoinsCommandFSFromFlag) {
		fmt.Println("Address %s is not valid", *sendCoinsCommandFSFromFlag)
		return
	}
	if len(*sendCoinsCommandFSToFlag) == 0 {
		fmt.Printf("Receiver address is not entered\n")
		return
	}
	if !checkWalletAddress(*sendCoinsCommandFSToFlag) {
		fmt.Println("Address %s is not valid", *sendCoinsCommandFSToFlag)
		return
	}
	if *sendCoinsCommandFSAmountFlag < 0 {
		fmt.Println("Amount of coins is not valid")
		return
	}

	bc, err := blockchain.NewBlockchain(*sendCoinsCommandFSDBAddress, []byte{})
	panicOnError(err)
	defer bc.Finalize()

	tx, err := blockchain.NewTransaction([]byte(*sendCoinsCommandFSFromFlag),
		[]byte(*sendCoinsCommandFSToFlag), *sendCoinsCommandFSAmountFlag, bc, ws)
	panicOnError(err)
	hash, err := bc.MineBlock([]*pb.Transaction{tx})
	panicOnError(err)
	fmt.Printf("Block mined, hash: %x\n", hash)
}

func (cli *CLIAppplication) createCommand() {
	createBlockchainFS := flag.NewFlagSet(createBlockchainCommand, flag.ExitOnError)
	createBlockchainFSAddressFlag := createBlockchainFS.String(createBlockchainCommandAddressFlag, "", "address of the blockchain founder")
	createBlockchainFSForceRecreateFlag := createBlockchainFS.Bool(createBlockchainCommandForceRecreateFlag, false, "force recreation of database base")
	createBlockchainDBAddress := createBlockchainFS.String(DBAddressFlag, defaultDbPath, "provides path to a specific blockchain db")

	err := createBlockchainFS.Parse(os.Args[2:])
	if err != nil {
		cli.printUsage()
		panic(err)
	}
	if !checkIsExt(*createBlockchainDBAddress, ".db") {
		fmt.Printf("%s is not a valid database\n", *createBlockchainDBAddress)
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
	if len(*createBlockchainFSAddressFlag) == 0 {
		fmt.Println("Founder address should not be empty")
		return
	}
	if !checkWalletAddress(*createBlockchainFSAddressFlag) {
		fmt.Println("Address %s is not valid", *createBlockchainFSAddressFlag)
		return
	}

	bc, err := blockchain.NewBlockchain(*createBlockchainDBAddress, []byte(*createBlockchainFSAddressFlag))
	panicOnError(err)
	defer bc.Finalize()

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
	if !checkIsExt(*deleteBlockchainDBAddress, ".db") {
		fmt.Printf("%s is not a valid database\n", *deleteBlockchainDBAddress)
		return
	}
	if removeDBFile(*deleteBlockchainDBAddress) {
		fmt.Printf("Database %s removed succesfully!\n", *deleteBlockchainDBAddress)
	}
}
