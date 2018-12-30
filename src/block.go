package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Version    uint64
	PrevHash   []byte
	MerkelRoot []byte
	Timestamp  uint64
	Difficulty uint64
	Nonce      uint64
	Hash       []byte
	Data       []byte
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

func (this *Block) Print() {
	fmt.Printf("version:%d\n", this.Version)
	fmt.Printf("prevHash:%x\n", this.PrevHash)
	fmt.Printf("merkelRoot:%x\n", this.MerkelRoot)
	fmt.Printf("timestamp:%d\n", this.Timestamp)
	fmt.Printf("difficulty:%d\n", this.Difficulty)
	fmt.Printf("nonce:%d\n", this.Nonce)
	fmt.Printf("hash:%x\n", this.Hash)
	fmt.Printf("data:%s\n", this.Data)
}

func NewBlock(prevHash []byte, data []byte) *Block {
	block := Block{
		Version:    1,
		PrevHash:   prevHash,
		MerkelRoot: []byte{},
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: 1,
		Nonce:      0,
		Hash:       []byte{},
		Data:       data,
	}
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
