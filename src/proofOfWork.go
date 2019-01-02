package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"
)

type ProofOfWork struct {
	_block  *Block
	_target *big.Int
}

func (this *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	block := this._block
	for {
		tmp := [][]byte{
			Uint64toSliceByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64toSliceByte(block.Timestamp),
			Uint64toSliceByte(block.Difficulty),
			Uint64toSliceByte(nonce),
		}
		blockInfo := bytes.Join(tmp, []byte{})
		hash = sha256.Sum256(blockInfo)
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])
		if tmpInt.Cmp(this._target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce
}

func NewProofOfWork(block *Block) *ProofOfWork {
	targetStr := "0010000000000000000000000000000000000000000000000000000000000000"
	target := big.Int{}
	target.SetString(targetStr, 16)
	pow := ProofOfWork{
		_block:  block,
		_target: &target,
	}
	return &pow
}

func Uint64toSliceByte(input uint64) []byte {
	buffer := bytes.Buffer{}
	err := binary.Write(&buffer, binary.BigEndian, input)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
