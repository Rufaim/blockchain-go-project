# Blockchain: go project

This project is a simplified blockchain on golang.
Mostly is follows [bitcoin protocols](https://en.bitcoin.it/wiki/Main_Page)


## How to use
### Before-the-launch
1. install a protobuf compiler for golang
2. run `make build_proto` command from the root of the repo to build protofiles
3. run `make run_tests` to test the package 
4. build the app with `go build -o bchain  *.go`

Now you should have executable named *bchain* inside the repository root

### Wallets
At first you require a couple of wallets.
You can create wallet with `./bchain wallet -new`.
Please note that wallet creation is stochastic, and addresses will be differend.

```
~$ ./bchain wallet -new`
New address: 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM
```

Running `./bchain wallet` will show you all generated wallet addresses.
Those addresses are fully valid Bitcoin addresses.
```
~$ ./bchain wallet
1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF
17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM
```

### Blockchain
Now let's create a blockchain itself.
```
~$ ./bchain create -address 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM
```

A wallet address provided during creation should now have coins for mining a genesis block.
To check the balance, you can use `./bchain balance` command and verify that creator has mining subsidy while the other ono's balance is zero.

```
~$ ./bchain balance -of 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM
Balance for 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM is 10
~$ ./bchain balance -of 1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF
Balance for 1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF is 0
```

Let's transfer some coins, shall we? Command `./bchain send` is used for that.
```
~$ ./bchain send -from 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM -to 1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF -amount 4
Block mined, hash: 00001a24a471471362f9c475d072b26b706ad136ebf41182a0580c5a5c7c6268
```
Current version of the chain does not implement a memory pool, therefore block is mined on each new transaction.
After the transfer you can see that the second wallet now have coins.
Note that in current version block is assumed to be mined by a sender, so the sender receives a mining subsidy.
```
~$ ./bchain balance -of 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM
Balance for 17Qsc97qvH5DkkvEjmZe72SUbUyyzZbHRM is 16
~$ ./bchain balance -of 1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF
Balance for 1Hn6QFFZ6hBLJWZzfvXRDAEiA3k2ZheGwF is 4
```
Run `./bchain show` to print the chain in a human-readable format

