package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"time"
	"unsafe"
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

func Uint64toSliceByte(input uint64) []byte {
	buffer := bytes.Buffer{}
	err := binary.Write(&buffer, binary.BigEndian, input)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func (this *Block) SetHash() {
	tmp := [][]byte{
		Uint64toSliceByte(this._version),
		this._prevHash,
		this._merkelRoot,
		Uint64toSliceByte(this._timestamp),
		Uint64toSliceByte(this._difficulty),
		Uint64toSliceByte(this._nonce),
		this._data,
	}
	blockInfo := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockInfo)
	this._hash = hash[:]
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
	fmt.Println(unsafe.Sizeof(*this))
}

func NewBlock(prevHash []byte, data []byte) *Block {
	block := Block{
		_version:    0,
		_prevHash:   prevHash,
		_merkelRoot: []byte{},
		_timestamp:  uint64(time.Now().Unix()),
		_difficulty: 1,
		_nonce:      0,
		_hash:       []byte{},
		_data:       data,
	}
	block.SetHash()
	return &block
}
