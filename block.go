package main

import (
	"crypto/sha256"
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
