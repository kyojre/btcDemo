package main

import (
	"fmt"
)

func (this *Block) Print() {
	fmt.Printf("version:%d\n", this.Version)
	fmt.Printf("prevHash:%x\n", this.PrevHash)
	fmt.Printf("merkelRoot:%x\n", this.MerkelRoot)
	fmt.Printf("timestamp:%d\n", this.Timestamp)
	fmt.Printf("difficulty:%d\n", this.Difficulty)
	fmt.Printf("nonce:%d\n", this.Nonce)
	fmt.Printf("hash:%x\n", this.Hash)
	fmt.Printf("Transactions:\n")
	for index, transaction := range this.Transactions {
		fmt.Printf("	index:%d\n", index)
		fmt.Printf("	txid:%x\n", transaction.TXID)
		fmt.Printf("	TXInputs\n")
		for inputIndex, input := range transaction.TXInputs {
			fmt.Printf("		inputIndex:%d\n", inputIndex)
			fmt.Printf("		txid:%x\n", input.TXID)
			fmt.Printf("		index:%d\n", input.Index)
			fmt.Printf("		sig:%s\n", input.Sig)
			fmt.Printf("		----\n")
		}
		fmt.Printf("	TXOutputs\n")
		for outputIndex, output := range transaction.TXOutputs {
			fmt.Printf("		outputIndex:%d\n", outputIndex)
			fmt.Printf("		value:%f\n", output.Value)
			fmt.Printf("		pubKeyHash:%s\n", output.PubKeyHash)
			fmt.Printf("		----\n")
		}
		fmt.Printf("	----\n")
	}
	fmt.Printf("----\n")
}

func (this *CLI) PrintBlockChain() {
	for blockChainIterator := this._blockChain.Iterator(); blockChainIterator.HasNext(); {
		block := blockChainIterator.Next()
		block.Print()
	}
}

func (this *CLI) GetBalance(address string) {
	utxos := this._blockChain.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("balance:%f\n", total)
}

func (this *CLI) Send(from string, to string, amount float64, miner string, data string) {
	transaction := NewTransaction(from, to, amount, this._blockChain)
	if transaction == nil {
		fmt.Printf("no_transaction\n")
		return
	}
	coinbaseTX := NewCoinbaseTX(miner, data)
	var transactions []*Transaction
	transactions = append(transactions, coinbaseTX)
	transactions = append(transactions, transaction)
	this._blockChain.AddBlock(transactions)
}

func (this *CLI) PrintWallets() {
	wallets := NewWallets()
	for address, wallet := range wallets.WalletsMap {
		fmt.Printf("address:%s\n", address)
		fmt.Printf("priKey:%v\n", wallet.PriKey)
		fmt.Printf("pubKey:%x\n", wallet.PubKey)
		fmt.Printf("----\n")
	}
}

func (this *CLI) CreateWallet() {
	wallets := NewWallets()
	wallets.CreateWallet()
}
