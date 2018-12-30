package main

import ()

func main() {
	blockChain := NewBlockChain()
	cli := CLI{
		_blockChain: blockChain,
	}
	cli.Run()
}
