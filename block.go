package main

import (
	"fmt"
	"time"
)

type Block struct {
	_version    uint64
	_prevHash   []byte
	_merkelRoot []byte
	_timestamp  uint64
	_difficulty uint64
	_nonce      uint64
	_hash       []byte
	_data       []byte
}

func (this *Block) Print() {
	fmt.Printf("version:%d\n", this._version)
	fmt.Printf("prevHash:%x\n", this._prevHash)
	fmt.Printf("merkelRoot:%x\n", this._merkelRoot)
	fmt.Printf("timestamp:%d\n", this._timestamp)
	fmt.Printf("difficulty:%d\n", this._difficulty)
	fmt.Printf("nonce:%d\n", this._nonce)
	fmt.Printf("hash:%x\n", this._hash)
	fmt.Printf("data:%s\n", this._data)
}

func NewBlock(prevHash []byte, data []byte) *Block {
	block := Block{
		_version:    1,
		_prevHash:   prevHash,
		_merkelRoot: []byte{},
		_timestamp:  uint64(time.Now().Unix()),
		_difficulty: 1,
		_nonce:      0,
		_hash:       []byte{},
		_data:       data,
	}
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block._hash = hash
	block._nonce = nonce
	return &block
}
