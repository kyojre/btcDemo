package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
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
	this.PubKeyHash = GetPubKeyHashByAddress(address)
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
		if bytes.Equal(input.TXID, []byte{}) && input.Index == -1 {
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

	blockChain.SignTransaction(&transaction, fromWallet.PriKey)

	return &transaction
}

func (this *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	if this.IsCoinbase() {
		return
	}
	txCopy := this.TrimmedCopy()
	for i, input := range txCopy.TXInputs {
		prevTx := prevTXs[string(input.TXID)]
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		signData := txCopy.TXID

		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signData)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		this.TXInputs[i].Sig = signature
	}
}

func (this *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range this.TXInputs {
		inputs = append(inputs, TXInput{
			TXID:   input.TXID,
			Index:  input.Index,
			Sig:    nil,
			PubKey: nil,
		})
	}
	for _, output := range this.TXOutputs {
		outputs = append(outputs, TXOutput{
			Value:      output.Value,
			PubKeyHash: output.PubKeyHash,
		})
	}
	return Transaction{
		TXID:      this.TXID,
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
}

func (this *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	if this.IsCoinbase() {
		return true
	}
	txCopy := this.TrimmedCopy()
	for i, input := range this.TXInputs {
		prevTx := prevTXs[string(input.TXID)]
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		signData := txCopy.TXID

		signature := input.Sig
		sigLen := len(signature)
		r := big.Int{}
		s := big.Int{}
		r.SetBytes(signature[:sigLen/2])
		s.SetBytes(signature[sigLen/2:])

		pubKey := input.PubKey
		keyLen := len(pubKey)
		x := big.Int{}
		y := big.Int{}
		x.SetBytes(pubKey[:keyLen/2])
		y.SetBytes(pubKey[keyLen/2:])

		pubKeyOrigin := ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     &x,
			Y:     &y,
		}

		res := ecdsa.Verify(&pubKeyOrigin, signData, &r, &s)
		if !res {
			return false
		}
	}
	return true
}
