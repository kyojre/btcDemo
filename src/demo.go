package main

import ()

func main() {
	blockChain := NewBlockChain("demo", "first")
	cli := CLI{
		_blockChain: blockChain,
	}
	cli.Run()
}
