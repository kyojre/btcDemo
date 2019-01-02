package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Version      uint64
	PrevHash     []byte
	MerkelRoot   []byte
	Timestamp    uint64
	Difficulty   uint64
	Nonce        uint64
	Hash         []byte
	Transactions []*Transaction
}

func (this *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(*this)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func (this *Block) MakeMerkelRoot() {
	this.MerkelRoot = []byte{}
}

func NewBlock(prevHash []byte, Transactions []*Transaction) *Block {
	block := Block{
		Version:      1,
		PrevHash:     prevHash,
		Timestamp:    uint64(time.Now().Unix()),
		Difficulty:   1,
		Nonce:        0,
		Hash:         []byte{},
		Transactions: Transactions,
	}
	block.MakeMerkelRoot()

	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

func Deserialize(data []byte) *Block {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
