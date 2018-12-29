package main

import (
	"fmt"
)

func main() {
	pblockChain := NewBlockChain()
	pblockChain.AddBlock([]byte("第二个区块"))
	for index, pblock := range pblockChain._blocks {
		fmt.Printf("index:%d\n", index)
		pblock.Print()
		fmt.Printf("\n")
	}
}
