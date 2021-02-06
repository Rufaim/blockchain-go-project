package main

import "github.com/Rufaim/blockchain/cmd/cli"

func main() {
	cliApp := cli.NewCLIAppplication()
	cliApp.Run()
}
