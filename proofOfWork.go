package main

import (
	"math/big"
)

type ProofOfWork struct {
	_block  *Block
	_target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	targetStr := "0000000000000000"
	target := big.Int{}
	target.SetString(targetStr, 16)
	pow := ProofOfWork{
		_block:  block,
		_target: &target,
	}
	return &pow
}
