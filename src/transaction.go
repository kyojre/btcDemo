package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 50.0

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXID  []byte
	Index int64
	Sig   string
}

type TXOutput struct {
	Value      float64
	PubKeyHash string
}

func (this *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(*this)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	this.TXID = hash[:]
}

func (this *Transaction) IsCoinbase() bool {
	if len(this.TXInputs) == 1 {
		input := this.TXInputs[0]
		if bytes.Equal(input.TXID, []byte{}) && input.Index != -1 {
			return true
		}
	}
	return false
}

func NewCoinbaseTX(address string) *Transaction {
	//挖矿交易的特点：
	//1.只有一个input
	//2.无需引用交易ID
	//3.无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写，一般是矿池名字
	input := TXInput{
		TXID:  []byte{},
		Index: -1,
		Sig:   "helloWorld",
	}
	output := TXOutput{
		Value:      reward,
		PubKeyHash: address,
	}
	transaction := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	transaction.SetHash()
	return &transaction
}
