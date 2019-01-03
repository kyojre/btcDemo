package main

import ()

func main() {
	wallets := NewWallets()
	address := wallets.CreateWallet()
	blockChain := NewBlockChain(address, "first")
	cli := CLI{
		_blockChain: blockChain,
	}
	cli.Run()
}
