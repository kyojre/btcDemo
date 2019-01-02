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

func NewCoinbaseTX(address string, data string) *Transaction {
	//挖矿交易的特点：
	//1.只有一个input
	//2.无需引用交易ID
	//3.无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写，一般是矿池名字
	input := TXInput{
		TXID:  []byte{},
		Index: -1,
		Sig:   data,
	}
	output := TXOutput{
		Value:      reward,
		PubKeyHash: address, //temp
	}
	transaction := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	transaction.SetHash()
	return &transaction
}

func NewTransaction(from string, to string, amount float64, blockChain *BlockChain) *Transaction {
	utxos, resValue := blockChain.FindNeedUTXOs(from, amount)

	if resValue < amount {
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{
				TXID:  []byte(id),
				Index: i,
				Sig:   from, //temp
			}
			inputs = append(inputs, input)
		}
	}
	output := TXOutput{
		Value:      amount,
		PubKeyHash: to, //temp
	}
	outputs = append(outputs, output)
	if resValue > amount {
		outputs = append(outputs, TXOutput{
			Value:      resValue - amount,
			PubKeyHash: from, //temp
		})
	}
	transaction := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	transaction.SetHash()
	return &transaction
}
