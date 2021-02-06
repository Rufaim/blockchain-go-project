package cli

const (
	defaultDbPath = "blockchain.db"
	DBAddressFlag = "path"
)

//// show command
const (
	showBlockchainCommand = "show"
)

//// add data commnd
const (
	addDataCommand         = "add"
	addDataCommandDataFlag = "data"
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
