package main

import (
	"fmt"
	"os"
	"strconv"
)

const USAGE = `	printChain
	getBalance ADDRES
	send FROM TO AMOUNT MINER DATA
	printWallets
	createWallets
`

type CLI struct {
	_blockChain *BlockChain
}

func (this *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(USAGE)
		return
	}

	cmd := args[1]
	switch cmd {
	case "printChain":
		this.PrintBlockChain()
	case "getBalance":
		if len(args) < 3 {
			fmt.Printf(USAGE)
			return
		}
		address := args[2]
		this.GetBalance(address)
	case "send":
		if len(args) < 6 {
			fmt.Printf(USAGE)
			return
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		this.Send(from, to, amount, miner, data)
	case "printWallets":
		this.PrintWallets()
	case "createWallet":
		this.CreateWallet()
	default:
		fmt.Printf(USAGE)
	}
}
