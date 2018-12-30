package main

import (
	"fmt"
	"os"
)

const USAGE = `
	addBlock --data DATA	 "add data to blockChain"
	printChain		 "print all blockChain data"
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
	case "addBlock":
		if len(args) < 3 {
			fmt.Printf(USAGE)
			return
		}
		blockData := args[2]
		this.AddBlock(blockData)
	case "printChain":
		this.PrintBlockChain()
	default:
		fmt.Printf(USAGE)
	}
}
