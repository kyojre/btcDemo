package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	_db                 *bolt.DB
	_currentHashPointer []byte
}

func (this *BlockChain) NewIterator() *BlockChainIterator {
	blockChainIterator := BlockChainIterator{
		_db:                 this._db,
		_currentHashPointer: this._tail,
	}
	return &blockChainIterator
}

func (this *BlockChainIterator) Next() *Block {
	if bytes.Compare(this._currentHashPointer, []byte{}) == 0 {
		return nil
	}
	var buffer []byte
	this._db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_CHAIN_BUCKET))
		if bucket == nil {
			log.Panic("no_bucket")
		} else {
			buffer = bucket.Get(this._currentHashPointer)
			if buffer == nil {
				log.Panic("no_block")
			}
		}
		return nil
	})
	block := Deserialize(buffer)
	this._currentHashPointer = block.PrevHash
	return block
}
