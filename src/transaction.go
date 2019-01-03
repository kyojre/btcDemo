package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const reward = 50.0

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXID   []byte
	Index  int64
	Sig    []byte
	PubKey []byte
}

type TXOutput struct {
	Value      float64
	PubKeyHash []byte
}

func (this *TXOutput) Lock(address string) {
	addressByte := base58.Decode(address)
	addressLen := len(addressByte)
	this.PubKeyHash = addressByte[1 : addressLen-4]
}

func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
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

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
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
		TXID:   []byte{},
		Index:  -1,
		Sig:    []byte(data),
		PubKey: []byte{},
	}
	output := NewTXOutput(reward, address)
	transaction := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{*output},
	}
	transaction.SetHash()
	return &transaction
}

func NewTransaction(from string, to string, amount float64, blockChain *BlockChain) *Transaction {
	wallets := NewWallets()
	fromWallet := wallets.WalletsMap[from]
	if fromWallet == nil {
		fmt.Printf("no_from_wallet\n")
		return nil
	}
	fromPubKey := fromWallet.PubKey
	fromPubKeyHash := HashPubKey(fromPubKey)

	utxos, resValue := blockChain.FindNeedUTXOs(fromPubKeyHash, amount)

	if resValue < amount {
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{
				TXID:   []byte(id),
				Index:  i,
				Sig:    []byte(from), //temp
				PubKey: fromPubKey,
			}
			inputs = append(inputs, input)
		}
	}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)
	if resValue > amount {
		leftOutput := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, *leftOutput)
	}
	transaction := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	transaction.SetHash()
	return &transaction
}
