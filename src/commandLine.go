package main

import (
	"fmt"
)

func (this *CLI) AddBlock(data string) {
	this._blockChain.AddBlock([]byte(data))
}

func (this *CLI) PrintBlockChain() {
	for blockChainIterator := this._blockChain.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		block.Print()
		fmt.Println("=================")
	}
}
