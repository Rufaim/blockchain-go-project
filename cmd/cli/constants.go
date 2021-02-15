package cli

const (
	defaultDbPath = "blockchain.db"
	DBAddressFlag = "path"
)

//// show command
const (
	showBlockchainCommand = "show"
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
