package main

import (
	"fmt"
)

func main() {
	blockChain := NewBlockChain()
	blockChain.AddBlock([]byte("第二个区块"))
	blockChainIterator := blockChain.NewIterator()
	for {
		block := blockChainIterator.Next()
		if block == nil {
			break
		}
		block.Print()
		fmt.Println("=================")
	}
}
