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

func (this *BlockChain) AddBlock(transactions []*Transaction) {
	this._db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_CHAIN_BUCKET))
		if bucket == nil {
			log.Panic("no_bucket")
		} else {
			block := NewBlock(this._tail, transactions)
			hash := block.Hash
			bucket.Put(hash, block.Serialize())
			bucket.Put([]byte("tail"), hash)
			this._tail = hash
		}
		return nil
	})
}

func NewBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address)
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

func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address)
	return NewBlock([]byte{}, []*Transaction{coinbase})
}

func (this *BlockChain) FindUTXOs(address string) []TXOutput {
	var utxos []TXOutput
	spentOutputs := make(map[string][]int64)
	for blockChainIterator := this.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		for _, transaction := range block.Transactions {
		OUTPUT:
			for i, output := range transaction.TXOutputs {
				if spentOutputs[string(transaction.TXID)] != nil {
					for _, j := range spentOutputs[string(transaction.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				//temp
				if output.PubKeyHash == address {
					utxos = append(utxos, output)
				}
			}
			if !transaction.IsCoinbase() {
				for _, input := range transaction.TXInputs {
					//temp
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXID)]
						indexArray = append(indexArray, input.Index)
					}
				}
			}
		}

	}
	return utxos
}
