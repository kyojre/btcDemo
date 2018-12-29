package main

import (
	"fmt"
)

func main() {
	pblockChain := NewBlockChain()
	pblockChain.AddBlock([]byte("第二个区块"))
	for index, pblock := range pblockChain._blocks {
		fmt.Printf("index:%d\nprevHash:%x\nhash:%x\ndata:%s\n\n", index, pblock._prevHash, pblock._hash, pblock._data)
	}
}
