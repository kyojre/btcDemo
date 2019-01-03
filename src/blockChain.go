package main

import (
	"bytes"
	"crypto/ecdsa"
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
	for _, transaction := range transactions {
		res := this.VerifyTransaction(transaction)
		if !res {
			//log.Panic("verify_false")
			return
		}
	}

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

func NewBlockChain(address string, data string) *BlockChain {
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
			genesisBlock := GenesisBlock(address, data)
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

func GenesisBlock(address string, data string) *Block {
	coinbase := NewCoinbaseTX(address, data)
	return NewBlock([]byte{}, []*Transaction{coinbase})
}

func (this *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var utxos []TXOutput
	spentOutputs := make(map[string]map[int64]bool)
	for blockChainIterator := this.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		for _, transaction := range block.Transactions {
			for i, output := range transaction.TXOutputs {
				if spentOutputs[string(transaction.TXID)] != nil {
					if spentOutputs[string(transaction.TXID)][int64(i)] {
						continue
					}
				}
				if bytes.Equal(output.PubKeyHash, pubKeyHash) {
					utxos = append(utxos, output)
				}
			}
			if !transaction.IsCoinbase() {
				for _, input := range transaction.TXInputs {
					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
						if spentOutputs[string(input.TXID)] == nil {
							spentOutputs[string(input.TXID)] = make(map[int64]bool)
						}
						spentOutputs[string(input.TXID)][input.Index] = true
					}
				}
			}
		}

	}
	return utxos
}

func (this *BlockChain) FindNeedUTXOs(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	utxos := make(map[string][]int64)
	var total float64 = 0.0
	spentOutputs := make(map[string]map[int64]bool)
	for blockChainIterator := this.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		for _, transaction := range block.Transactions {
			for i, output := range transaction.TXOutputs {
				if spentOutputs[string(transaction.TXID)] != nil {
					if spentOutputs[string(transaction.TXID)][int64(i)] {
						continue
					}
				}
				if bytes.Equal(output.PubKeyHash, pubKeyHash) {
					utxos[string(transaction.TXID)] = append(utxos[string(transaction.TXID)], int64(i))
					total += output.Value
					if total >= amount {
						return utxos, total
					}
				}
			}
			if !transaction.IsCoinbase() {
				for _, input := range transaction.TXInputs {
					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
						if spentOutputs[string(input.TXID)] == nil {
							spentOutputs[string(input.TXID)] = make(map[int64]bool)
						}
						spentOutputs[string(input.TXID)][input.Index] = true
					}
				}
			}
		}

	}
	return utxos, total
}

func (this *BlockChain) FindUTXOTransactions(pubKeyHash []byte) []*Transaction {
	var txs []*Transaction
	spentOutputs := make(map[string]map[int64]bool)
	for blockChainIterator := this.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		for _, transaction := range block.Transactions {
			for i, output := range transaction.TXOutputs {
				if spentOutputs[string(transaction.TXID)] != nil {
					if spentOutputs[string(transaction.TXID)][int64(i)] {
						continue
					}
				}
				if bytes.Equal(output.PubKeyHash, pubKeyHash) {
					txs = append(txs, transaction)
				}
			}
			if !transaction.IsCoinbase() {
				for _, input := range transaction.TXInputs {
					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
						if spentOutputs[string(input.TXID)] == nil {
							spentOutputs[string(input.TXID)] = make(map[int64]bool)
						}
						spentOutputs[string(input.TXID)][input.Index] = true
					}
				}
			}
		}
	}
	return txs
}

func (this *BlockChain) FindTransactionByTXid(TXID []byte) *Transaction {
	for blockChainIterator := this.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		for _, transaction := range block.Transactions {
			if bytes.Equal(transaction.TXID, TXID) {
				return transaction
			}
		}
	}
	return nil
}

func (this *BlockChain) SignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) {
	prevTXs := make(map[string]*Transaction)
	for _, input := range transaction.TXInputs {
		tx := this.FindTransactionByTXid(input.TXID)
		if tx == nil {
			log.Panic("no_tx")
		}
		prevTXs[string(input.TXID)] = tx
	}
	transaction.Sign(privateKey, prevTXs)
}

func (this *BlockChain) VerifyTransaction(transaction *Transaction) bool {
	if transaction.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]*Transaction)
	for _, input := range transaction.TXInputs {
		tx := this.FindTransactionByTXid(input.TXID)
		if tx == nil {
			log.Panic("no_tx")
		}
		prevTXs[string(input.TXID)] = tx
	}
	return transaction.Verify(prevTXs)
}
