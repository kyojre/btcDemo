package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	_prevHash []byte
	_hash     []byte
	_data     []byte
}

func (this *Block) SetHash() {
	blockInfo := append(this._prevHash, this._data...)
	hash := sha256.Sum256(blockInfo)
	this._hash = hash[:]
}

func NewBlock(prevHash []byte, data []byte) *Block {
	block := Block{
		_prevHash: prevHash,
		_hash:     []byte{},
		_data:     data,
	}
	block.SetHash()
	return &block
}

type BlockChain struct {
	_blocks []*Block
}

func (this *BlockChain) AddBlock(data []byte) {
	plastBlock := this._blocks[len(this._blocks)-1]
	prevHash := plastBlock._hash
	pblock := NewBlock(prevHash, data)
	this._blocks = append(this._blocks, pblock)
}

func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	blockChain := BlockChain{
		_blocks: []*Block{genesisBlock},
	}
	return &blockChain
}

func GenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("创世区块"))
}

func main() {
	pblockChain := NewBlockChain()
	pblockChain.AddBlock([]byte("第二个区块"))
	for index, pblock := range pblockChain._blocks {
		fmt.Printf("index:%d\nprevHash:%x\nhash:%x\ndata:%s\n\n", index, pblock._prevHash, pblock._hash, pblock._data)
	}
}
