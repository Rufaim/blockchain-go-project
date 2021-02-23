package cli

const (
	defaultDbPath        = "blockchain.db"
	defaultWalletSetPath = "wallet_set.ws"
	DBAddressFlag        = "dbpath"
	WSAddressFlag        = "wspath"
)

//// show command
const (
	showBlockchainCommand = "show"
)

//// wallet
const (
	walletOperationsCommand        = "wallet"
	walletOperationsCommandNewFlag = "new"
)

//// balance command
const (
	balanceOfWalletCommand       = "balance"
	balanceOfWalletCommandOFFlag = "of"
)

//// send command
const (
	sendCoinsCommand           = "send"
	sendCoinsCommandFromFlag   = "from"
	sendCoinsCommandToFlag     = "to"
	sendCoinsCommandAmountFlag = "amount"
)

//// create command
const (
	createBlockchainCommand                  = "create"
	createBlockchainCommandForceRecreateFlag = "f"
	createBlockchainCommandAddressFlag       = "address"
)

//// delete command
const (
	deleteBlockchainCommand = "delete"
)

//// help command
const (
	helpBlockchainCommand      = "help"
	helpBlockchainCommandHFlag = "h"
)
