package main

import ()

const address = "demo"

func main() {
	blockChain := NewBlockChain(address)
	cli := CLI{
		_blockChain: blockChain,
	}
	cli.Run()
}
