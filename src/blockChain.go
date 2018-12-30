package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const BLOCK_CHAIN_DB = "../db/blockChainDB"
const BLOCK_CHAIN_BUCKET = "blockChainBucket"

type BlockChain struct {
	_db   *bolt.DB
	_tail []byte
}

func (this *BlockChain) AddBlock(data []byte) {
	this._db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_CHAIN_BUCKET))
		if bucket == nil {
			log.Panic("no_bucket")
		} else {
			block := NewBlock(this._tail, data)
			hash := block.Hash
			bucket.Put(hash, block.Serialize())
			bucket.Put([]byte("tail"), hash)
			this._tail = hash
		}
		return nil
	})
}

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(BLOCK_CHAIN_DB, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//defer db.Close()
	var tail []byte
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_CHAIN_BUCKET))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(BLOCK_CHAIN_BUCKET))
			if err != nil {
				log.Panic(err)
			}
			genesisBlock := GenesisBlock()
			tail = genesisBlock.Hash
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("tail"), tail)
		} else {
			tail = bucket.Get([]byte("tail"))
		}
		return nil
	})
	blockChain := BlockChain{
		_db:   db,
		_tail: tail,
	}
	return &blockChain
}

func GenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("genesisBlock"))
}
